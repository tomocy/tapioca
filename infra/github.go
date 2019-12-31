package infra

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/tomocy/tapioca/domain"
	infragithub "github.com/tomocy/tapioca/infra/github"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const authorizationRedirectPath = "/tapioca/authorization"

func NewGitHub() *GitHub {
	createWorkspace()
	return &GitHub{
		client: oauth2Client{
			conf: oauth2.Config{
				ClientID:     "5a24485cf2fe2ca8fab4",
				ClientSecret: "63a169863256d15eca02ac6ade415f93b2692e28",
				RedirectURL:  "http://localhost/tapioca/authorization",
				Scopes: []string{
					"read:repo", "read:user",
				},
				Endpoint: github.Endpoint,
			},
		},
	}
}

type GitHub struct {
	client oauth2Client
}

type oauth2Client struct {
	state string
	conf  oauth2.Config
}

func (g *GitHub) FetchRepos(_ string) ([]*domain.Repo, error) {
	var rs infragithub.Repos
	if err := g.fetch(
		fmt.Sprintf("https://api.github.com/user/repos"),
		nil,
		&rs,
	); err != nil {
		return nil, err
	}

	return rs.Adapt(), nil
}

func (g *GitHub) FetchCommits(owner, repo string, params domain.Params) (domain.Commits, error) {
	ids, err := g.fetchCommitIDs(owner, repo, params)
	if err != nil {
		return nil, err
	}

	cs := make(domain.Commits, len(ids))
	for i, id := range ids {
		c, err := g.FetchCommit(owner, repo, id)
		if err != nil {
			return nil, err
		}
		cs[i] = c
	}

	return cs, nil
}

func (g *GitHub) fetchCommitIDs(owner, repo string, params domain.Params) ([]string, error) {
	var cs infragithub.Commits
	if err := g.fetch(
		fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo),
		g.parseParams(params),
		&cs,
	); err != nil {
		return nil, err
	}

	ids := make([]string, len(cs))
	for i, c := range cs {
		ids[i] = c.SHA
	}

	return ids, nil
}

func (g *GitHub) parseParams(params domain.Params) url.Values {
	vs := make(url.Values)
	if params.Author != "" {
		vs.Set("author", params.Author)
	}
	if !params.Since.IsZero() {
		vs.Set("since", params.Since.Format(time.RFC3339))
	}
	if !params.Until.IsZero() {
		vs.Set("until", params.Until.Format(time.RFC3339))
	}

	return vs
}

func (g *GitHub) FetchCommit(owner, repo, id string) (*domain.Commit, error) {
	var c infragithub.Commit
	if err := g.fetch(
		fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, id),
		nil,
		&c,
	); err != nil {
		return nil, err
	}

	return c.Adapt(), nil
}

func (g *GitHub) fetch(dstURI string, params url.Values, dst interface{}) error {
	tok, err := g.retieveAuthorization()
	if err != nil {
		return err
	}

	if err := g.do(&oauthReq{
		tok:    tok,
		method: http.MethodGet,
		uri:    dstURI,
		params: params,
	}, dst); err != nil {
		return err
	}

	return g.saveConfig(oauth2Config{
		AccessToken: tok,
	})
}

func (g *GitHub) saveConfig(conf oauth2Config) error {
	if loaded, err := loadConfig(); err == nil {
		loaded.GitHub = conf
		return saveConfig(loaded)
	}

	return saveConfig(&config{
		GitHub: conf,
	})
}

func (g *GitHub) retieveAuthorization() (*oauth2.Token, error) {
	if conf, err := loadConfig(); err == nil {
		return conf.GitHub.AccessToken, nil
	}

	url := g.oauthCodeURL()
	fmt.Printf("open this link: %s\n", url)

	return g.handleAuthorizationRedirect()
}

func (g *GitHub) oauthCodeURL() string {
	g.setRandomState()
	return g.client.conf.AuthCodeURL(g.client.state)
}

func (g *GitHub) setRandomState() {
	g.client.state = fmt.Sprintf("%d", rand.Intn(10000))
}

func (g *GitHub) handleAuthorizationRedirect() (*oauth2.Token, error) {
	tokCh, errCh := g.handleAsyncAuthorizationRedirect()
	select {
	case tok := <-tokCh:
		return tok, nil
	case err := <-errCh:
		return nil, err
	}
}

func (g *GitHub) handleAsyncAuthorizationRedirect() (<-chan *oauth2.Token, <-chan error) {
	tokCh, errCh := make(chan *oauth2.Token), make(chan error)
	go func() {
		defer func() {
			close(tokCh)
			close(errCh)
		}()

		http.Handle(authorizationRedirectPath, g.handlerForAuthorizationRedirect(tokCh, errCh))
		if err := http.ListenAndServe(":80", nil); err != nil {
			errCh <- err
			return
		}
	}()

	return tokCh, errCh
}

func (g *GitHub) handlerForAuthorizationRedirect(tokCh chan<- *oauth2.Token, errCh chan<- error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		state, code := q.Get("state"), q.Get("code")
		if err := g.checkState(state); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errCh <- err
			return
		}

		tok, err := g.client.conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errCh <- err
			return
		}

		tokCh <- tok
	})
}

func (g *GitHub) checkState(state string) error {
	stored := g.client.state
	g.client.state = ""
	if state != stored {
		return errors.New("invalid state")
	}

	return nil
}

func (g *GitHub) do(r *oauthReq, dst interface{}) error {
	resp, err := r.do(g.client.conf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return readJSON(resp.Body, dst)
}

type oauthReq struct {
	tok    *oauth2.Token
	method string
	uri    string
	params url.Values
}

func (r *oauthReq) do(conf oauth2.Config) (*http.Response, error) {
	client := conf.Client(oauth2.NoContext, r.tok)
	if r.method != http.MethodGet {
		return client.PostForm(r.uri, r.params)
	}

	parsed, err := url.Parse(r.uri)
	if err != nil {
		return nil, err
	}
	if r.params != nil {
		parsed.RawQuery = r.params.Encode()
	}

	return client.Get(parsed.String())
}
