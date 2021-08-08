import axios, { AxiosRequestConfig, AxiosResponse, AxiosError } from "axios";

const LOGIN_FORM_ID: string = "login";
const LOGIN_FORM_USERNAME_ID: string = "username";
const LOGIN_FORM_PASSWORD_ID: string = "password";
const LOGIN_FORM_ALERT_ID: string = "login-form-alert";
const LOGIN_API_PATH: string = "http://localhost:2101/api/v1/login";

let headers = {
  Accept: "application/json",
};

axios
  .get(LOGIN_API_PATH, { headers })
  .then((response) => {
    console.log(response);
    const loginID = response.data["login_id"];
    console.log(loginID);
  })
  .catch((error: AxiosError) => {
    const errorMessage = getErrorMessageOrDefault(error);
    const dangerAlertHTML = makeDangerAlert(errorMessage);
    document.getElementById(LOGIN_FORM_ALERT_ID).innerHTML = dangerAlertHTML;
  });

document.getElementById(LOGIN_FORM_ID).addEventListener("submit", (e) => {
  e.preventDefault();
  postLogin();
});

const postLogin = () => {
  const username: string = (
    document.getElementById(LOGIN_FORM_USERNAME_ID) as HTMLInputElement
  ).value;
  const password: string = (
    document.getElementById(LOGIN_FORM_PASSWORD_ID) as HTMLInputElement
  ).value;
  axios
    .put(LOGIN_API_PATH, { username, password }, { withCredentials: true })
    .then((_: AxiosResponse) => {
      document.getElementById(LOGIN_FORM_ALERT_ID).innerHTML = "";
    })
    .catch((error: AxiosError) => {
      const errorMessage = getErrorMessageOrDefault(error);
      const dangerAlertHTML = makeDangerAlert(errorMessage);
      document.getElementById(LOGIN_FORM_ALERT_ID).innerHTML = dangerAlertHTML;
    });
};

const makeDangerAlert = (alertMessage: string) => {
  return `<div class="alert alert-danger" role="alert">${alertMessage}</div>`;
};

const defaultLoginError: string =
  "Sorry! Something has gone wrong. Please try again.";

const getErrorMessageOrDefault = (error: AxiosError) => {
  let errorMessage = error.response.data.errorMessage;
  if (!errorMessage) {
    errorMessage = defaultLoginError;
  }
  return errorMessage;
};
