package handlers

import (
	"net/http"
)

// redirectWithError redirects back to main page with error message
func redirectWithError(w http.ResponseWriter, r *http.Request, message string) {
	// redirect back to main page with error as query parameter
	http.Redirect(w, r, "/?error="+message, http.StatusSeeOther)
}

// redirectWithMessage redirects back to main page with success message
func redirectWithMessage(w http.ResponseWriter, r *http.Request, message string) {
	// redirect back to main page with success as query parameter
	http.Redirect(w, r, "/?success="+message, http.StatusSeeOther)
}
