
const LOGIN_FORM_USERNAME_ID = "username";
const LOGIN_FORM_PASSWORD_ID = "password";
const LOGIN_PATH = "/login";

function postLogin() {
    let username = document.getElementById(LOGIN_FORM_USERNAME_ID).value;
    let password = document.getElementById(LOGIN_FORM_PASSWORD_ID).value;
    axios.post(LOGIN_PATH, {
        username: username,
        password: password
    })
        .then(function (response) {
            console.log(response);
        })
        .catch(function (error) {
            console.log(error);
        });
}

