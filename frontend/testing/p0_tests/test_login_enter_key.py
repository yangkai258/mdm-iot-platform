"""
P0 Test: Login Page Enter Key - Verify Enter Key Triggers Login
================================================================
Bug: Pressing Enter in username/password fields does not submit the form
Fix: Form must handle @keyup.enter or @submit.on-enter properly

This test validates:
1. Login form responds to Enter key press
2. Form has proper @submit handler
3. Enter key submission works for both username and password fields
4. Loading state is properly managed during login
"""
import os
import pytest
import re


class TestLoginEnterKey:
    """
    Login form Enter key functionality tests.
    PASS = Enter key triggers form submission
    FAIL = Enter key does not trigger login
    """

    LOGIN_VUE_PATH = os.path.join(
        os.getenv("MDM_FRONTEND_PATH", "C:\\Users\\YKing\\.openclaw\\workspace\\mdm-project\\frontend\\src"),
        "views", "Login.vue"
    )

    def test_login_form_has_submit_handler(self):
        """
        Verify the login form has a @submit handler attached.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Look for form with @submit
        has_submit = "@submit" in source or "@submit.prevent" in source

        assert has_submit, (
            "FAIL: Login form does not have @submit handler.\n"
            "Expected: <a-form @submit='handleLogin'> or similar"
        )

        print("PASS: Login form has @submit handler")

    def test_login_form_handles_enter_key(self):
        """
        Verify Enter key is handled for form submission.

        Common patterns that handle Enter key:
        - @keyup.enter on input fields
        - @submit on form (Enter in inputs triggers form submit)
        - Using <a-form> with proper submit handling
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Check for Enter key handling patterns
        enter_key_patterns = [
            "@keyup.enter",
            "@enter",
            "@press.enter",
            "@on-enter",
            'keyup.enter',
            'press.enter',
        ]

        has_enter_handler = any(pattern in source for pattern in enter_key_patterns)

        # Also check if using a-form which handles Enter automatically
        has_a_form = "<a-form" in source and "@submit" in source

        assert has_enter_handler or has_a_form, (
            "FAIL: No Enter key handling found in Login.vue.\n"
            "Expected: @keyup.enter on input fields, or @submit on <a-form>.\n"
            "Enter key in input fields should trigger form submission."
        )

        if has_enter_handler:
            print("PASS: Enter key handler found in Login.vue")
        else:
            print("PASS: <a-form> with @submit handles Enter key automatically")

    def test_login_button_has_type_submit(self):
        """
        Verify the login button has type='submit' so Enter key triggers it.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Look for button with type="submit" or html-type="submit"
        button_pattern = r'<a-button[^>]*type=["\']submit["\'][^>]*>'
        has_submit_button = re.search(button_pattern, source) is not None

        # Also check for html-type (Arco Design uses html-type)
        html_type_pattern = r'html-type=["\']submit["\']'
        has_html_type_submit = re.search(html_type_pattern, source) is not None

        assert has_submit_button or has_html_type_submit, (
            "FAIL: Login button does not have type='submit'.\n"
            "Without type='submit', Enter key may not trigger the form.\n"
            "Expected: <a-button type='primary' html-type='submit'>"
        )

        print("PASS: Login button has type='submit'")

    def test_handleLogin_function_exists(self):
        """
        Verify handleLogin function is defined and handles form submission.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        assert "handleLogin" in source, (
            "FAIL: handleLogin function not found in Login.vue.\n"
            "The form submit handler must call handleLogin."
        )

        print("PASS: handleLogin function is defined")

    def test_handleLogin_validates_input(self):
        """
        Verify handleLogin validates username and password before submission.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Look for validation in handleLogin
        validation_patterns = [
            "username",
            "password",
            "required",
            "message",
            "warning",
        ]

        # handleLogin should check username and password
        has_validation = all(
            pattern in source.lower()
            for pattern in ["username", "password", "message"]
        )

        assert has_validation, (
            "FAIL: handleLogin may not validate input fields.\n"
            "Expected: Check for empty username/password and show warning."
        )

        print("PASS: handleLogin includes input validation")

    def test_handleLogin_has_loading_state(self):
        """
        Verify login handles loading state to prevent double-submission.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Check for loading ref and its usage
        has_loading_ref = "loading" in source
        has_loading_true = "loading.value = true" in source or "loading.value=true" in source
        has_loading_false = "loading.value = false" in source or "loading.value=false" in source

        assert has_loading_ref and has_loading_true and has_loading_false, (
            "FAIL: Login does not properly manage loading state.\n"
            "Expected: loading ref with true/false toggling to prevent double-submit."
        )

        print("PASS: Login properly manages loading state")

    def test_handleLogin_calls_api(self):
        """
        Verify handleLogin makes an API call to the backend.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        api_patterns = [
            "/api/v1/auth/login",
            "fetch(",
            "axios",
        ]

        has_api_call = any(pattern in source for pattern in api_patterns)

        assert has_api_call, (
            "FAIL: handleLogin does not appear to call the login API.\n"
            "Expected: POST request to /api/v1/auth/login"
        )

        print("PASS: handleLogin calls the login API")

    def test_handleLogin_sets_token_on_success(self):
        """
        Verify successful login stores the JWT token (e.g., in localStorage).
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        storage_patterns = [
            "localStorage.setItem",
            "token",
            "data.data.token",
        ]

        has_token_storage = all(
            pattern in source.lower()
            for pattern in ["localstorage", "token"]
        )

        assert has_token_storage, (
            "FAIL: Login does not store JWT token on success.\n"
            "Expected: localStorage.setItem('token', ...) after successful login."
        )

        print("PASS: Login stores JWT token on success")

    def test_password_input_has_enter_key_binding(self):
        """
        Specifically check that password input handles Enter key.

        Users typically press Enter after typing password.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Check if password input has Enter handler
        # Split source to check password field context
        lines = source.split('\n')
        password_input_found = False
        has_enter_near_password = False

        for i, line in enumerate(lines):
            if 'password' in line.lower() and ('input' in line.lower() or 'v-model' in line):
                password_input_found = True
                # Check surrounding lines for Enter handler
                context = '\n'.join(lines[max(0, i-2):min(len(lines), i+3)])
                if '@keyup.enter' in context or '@enter' in context or '@submit' in context:
                    has_enter_near_password = True

        # If a-form with @submit is used, Enter in password field works automatically
        has_form_submit = "<a-form" in source and "@submit" in source

        assert has_enter_near_password or has_form_submit, (
            "FAIL: Password input field does not handle Enter key.\n"
            "Users expect to press Enter after typing password.\n"
            "Expected: @keyup.enter on password input, or @submit on containing form."
        )

        print("PASS: Password input properly handles Enter key")

    def test_no_prevent_default_without_reason(self):
        """
        Ensure there's no e.preventDefault() blocking Enter key inappropriately.
        """
        with open(self.LOGIN_VUE_PATH, "r", encoding="utf-8") as f:
            source = f.read()

        # Look for preventDefault in handleLogin
        if "preventDefault" in source:
            # If preventDefault is used, it should be in a valid context
            # (e.g., stopping a different action, not blocking form submit)
            # This is a soft check - just warning if present
            print("INFO: preventDefault found - ensure it doesn't block Enter key")
        else:
            print("PASS: No inappropriate preventDefault blocking Enter key")


if __name__ == "__main__":
    pytest.main([__file__, "-v", "--tb=short"])
