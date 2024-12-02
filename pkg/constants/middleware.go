package constants

import "github.com/go-playground/validator/v10"

type ContextKey string

const (
	UserKey        ContextKey = "user"
	SessionKey     ContextKey = "session"
	NavItemsKey    ContextKey = "navItems"
	AllNavItemsKey ContextKey = "allNavItems"
	TxKey          ContextKey = "tx"
	ParamsKey      ContextKey = "params"
	LoggerKey      ContextKey = "logger"
	LogoKey        ContextKey = "logo"
	HeadKey        ContextKey = "head"
	TabsKey        ContextKey = "tabs"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())
