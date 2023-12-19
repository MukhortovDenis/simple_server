package simple_server

import (
	"encoding/json"
	"fmt"

	"vsu/config"
	"vsu/internal/auth/cache"
	"vsu/internal/auth/model"

	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	accCache  *cache.AccountCache
	permCache *cache.PermissionsCache
}

type RequestRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RequestAuthorize struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseAuthorize struct {
	Login       string                 `json:"login"`
	Permissions []model.PermissionType `json:"permissions"`
}

func NewService(accCache *cache.AccountCache, permCache *cache.PermissionsCache) *Service {
	return &Service{
		accCache:  accCache,
		permCache: permCache,
	}
}

func (s *Service) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/api/v1/register":
		if ctx.IsPost() {
			s.register(ctx)
		} else {
			ctx.Error("Unsupported method", fasthttp.StatusNotFound)
		}
	case "/api/v1/authorize":
		if ctx.IsGet() {
			s.authorize(ctx)
		} else {
			ctx.Error("Unsupported method", fasthttp.StatusNotFound)
		}
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func (s *Service) Start(cfg *config.Env) error {
	if err := s.addAdminAccounts(cfg); err != nil {
		return fmt.Errorf("admin account: %w", err)
	}

	fmt.Println("Waiting some requests)))")

	if err := fasthttp.ListenAndServe(cfg.HTTPListenAddr, s.HandleFastHTTP); err != nil {
		return fmt.Errorf("fasthttp error: %w", err)
	}

	return nil
}

func (s *Service) register(ctx *fasthttp.RequestCtx) {
	var req RequestRegister
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)

		return
	}

	if req.Login == "" || req.Password == "" {
		ctx.Error("Empty login or password", fasthttp.StatusBadRequest)

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		ctx.Error("Empty login or password", fasthttp.StatusInternalServerError)

		return
	}

	if err = s.accCache.Set(req.Login, string(hash)); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

		return
	}

	if err = s.permCache.Set(req.Login, model.DefaultUser); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

		return
	}
}

func (s *Service) authorize(ctx *fasthttp.RequestCtx) {
	var req RequestAuthorize
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)

		return
	}

	if req.Login == "" || req.Password == "" {
		ctx.Error("Empty login or password", fasthttp.StatusBadRequest)

		return
	}

	hash, err := s.accCache.GetPass(req.Login)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)

		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)

		return
	}

	perms, err := s.permCache.GetPermissions(req.Login)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

		return
	}

	resp := ResponseAuthorize{
		Login:       req.Login,
		Permissions: perms,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)

		return
	}

	ctx.Response.SetBody(body)
}

func (s *Service) addAdminAccounts(cfg *config.Env) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err = s.accCache.Set(cfg.AdminLogin, string(hash)); err != nil {
		return err
	}
	if err = s.permCache.Set(cfg.AdminLogin, model.AdminUser); err != nil {
		return err
	}

	if err = s.permCache.Set(cfg.AdminLogin, model.DefaultUser); err != nil {
		return err
	}

	return nil
}
