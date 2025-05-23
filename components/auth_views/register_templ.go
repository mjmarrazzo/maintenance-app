// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package auth_views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/mjmarrazzo/maintenance-app/components/common"
	"github.com/mjmarrazzo/maintenance-app/components/common/form"
)

type RegisterProps struct {
	Email string
}

func Register(props RegisterProps) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"hero min-h-screen\" style=\"background-image: url(/public/illustrations/hero.png);\"><div class=\"hero-overlay bg-opacity-60\"></div><div class=\"flex justify-center items-center\"><div id=\"registration-card\" class=\"card backdrop-blur-lg bg-white/20 shadow-xl \"><div class=\"card-body\"><h2 class=\"text-3xl font-bold text-white text-center\">Welcome to Groundwork!</h2><form id=\"registration-form\" class=\"flex flex-col gap-4 p-4\" hx-post=\"/register\" hx-indicator=\"#form-spinner\" hx-disabled-elt=\"button[type=submit]\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Input(form.InputProps{
				ID:           "email",
				Label:        "Email",
				Type:         "email",
				Value:        props.Email,
				Hint:         "Enter a valid email address",
				Autocomplete: form.AutocompleteUsername,
				IsRequired:   true,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Password(form.PasswordProps{
				ID:           "password",
				Label:        "Password",
				Autocomplete: form.AutocompleteCurrentPassword,
				Hint:         "Required",
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Password(form.PasswordProps{
				ID:           "password-confirmation",
				Label:        "Confirm Password",
				Autocomplete: form.AutocompleteCurrentPassword,
				Hint:         "Passwords must match",
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Input(form.InputProps{
				ID:           "first_name",
				Label:        "First Name",
				Type:         "text",
				Hint:         "Required",
				Autocomplete: form.AutocompleteGivenName,
				IsRequired:   true,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = form.Input(form.InputProps{
				ID:           "last_name",
				Label:        "Last Name",
				Type:         "text",
				Hint:         "Required",
				Autocomplete: form.AutocompleteFamilyName,
				IsRequired:   true,
			}).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<button class=\"btn btn-primary btn-block mt-2 shadow-md\" type=\"submit\">Login</button></form></div></div></div></div><script>\n            // Function to check passwords match on dirty of password confirmation\n            const passwordInput = document.getElementById('password');\n            const passwordConfirmationInput = document.getElementById('password-confirmation');\n            const passwordHint = passwordConfirmationInput.nextElementSibling;\n            passwordConfirmationInput.addEventListener('input', function(event) {\n                if (passwordInput.value !== passwordConfirmationInput.value) {\n                    passwordConfirmationInput.setCustomValidity(\"Passwords do not match\");\n                } else {\n                    passwordConfirmationInput.setCustomValidity(\"\");\n                }\n            });\n\n        </script>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = common.BaseHtml("Register").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
