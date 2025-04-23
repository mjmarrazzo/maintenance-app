package form

type autocomplete string

const (
	AutocompleteNewPassword     autocomplete = "new-password"
	AutocompleteCurrentPassword autocomplete = "current-password"
	AutocompleteUsername        autocomplete = "username"
	AutocompleteGivenName       autocomplete = "given-name"
	AutocompleteFamilyName      autocomplete = "family-name"
)
