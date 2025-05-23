// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package category_views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/mjmarrazzo/maintenance-app/components/common/form"
	"github.com/mjmarrazzo/maintenance-app/domain"
)

type FormProps struct {
	IsEdit   bool
	Category *domain.Category
}

func Form(props FormProps) templ.Component {
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<h3 class=\"text-lg font-bold\" id=\"dialog-title\" hx-swap-oob=\"#dialog-title\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.IsEdit {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "Edit Category")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "Create Category")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</h3><div class=\"p-4\"><form")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if props.IsEdit {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, " hx-put=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var2 string
			templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/categories/%d", props.Category.ID))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `components/category_views/form.templ`, Line: 25, Col: 61}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, " hx-post=\"/categories\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, " hx-target=\"#category-modal-content\" hx-swap=\"outerHTML\" hx-indicator=\"#form-spinner\" hx-disabled-elt=\".modal-action button\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.Input(form.InputProps{
			ID:         "name",
			Label:      "Name",
			Value:      safeCategory(props.Category).Name,
			Type:       "text",
			IsRequired: true,
			Hint:       "Required",
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = form.TextArea(form.TextAreaProps{
			ID:         "description",
			Label:      "Description",
			Value:      safeCategory(props.Category).Description,
			Rows:       3,
			IsRequired: false,
		}).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "<div class=\"modal-action\"><button type=\"button\" class=\"btn\" onclick=\"category_modal.close()\">Cancel</button> <button type=\"submit\" class=\"btn btn-primary\">Save Changes <span id=\"form-spinner\" class=\"htmx-indicator\"><span class=\"loading loading-spinner loading-md\"></span></span></button></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func safeCategory(category *domain.Category) *domain.Category {
	if category == nil {
		return &domain.Category{}
	}
	return category
}

var _ = templruntime.GeneratedTemplate
