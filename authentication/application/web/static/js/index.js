
const LOGIN_FORM_USERNAME_ID = "username";
const LOGIN_FORM_PASSWORD_ID = "password";
const LOGIN_FORM_ALERT_ID = "login-form-alert";
const LOGIN_API_PATH = "/api/v1/login";

const REGISTER_FORM_CONFIRM_PASSWORD_ID = "password";
const REGISTER_FORM_PASSWORD_ID = "confirm-password";
const USERS_API_PATH = "/api/v1/users";


function postLogin() {
    let username = document.getElementById(LOGIN_FORM_USERNAME_ID).value;
    let password = document.getElementById(LOGIN_FORM_PASSWORD_ID).value;
    axios.post(LOGIN_API_PATH, {
        username: username,
        password: password
    })
        .then(function (response) {
            document.getElementById(LOGIN_FORM_ALERT_ID).innerHTML = "";
            console.log(response);
        })
        .catch(function (error) {
            let errorMessage = getErrorMessageOrDefault(error);
            let dangerAlertHTML = makeDangerAlert(errorMessage);
            document.getElementById(LOGIN_FORM_ALERT_ID).innerHTML = dangerAlertHTML;
        });
}


function makeDangerAlert(alertMessage) {
    return `<div class="alert alert-danger" role="alert">${alertMessage}</div>`;
}


const defaultLoginError = "Sorry! Something has gone wrong. Please try again.";


function getErrorMessageOrDefault(error) {
    let error_message = error.response.data.errorMessage;
    if (!error_message) {
        error_message = defaultLoginError;
    }
    return error_message
}