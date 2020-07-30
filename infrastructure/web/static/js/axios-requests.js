
function postLogin() {
    let loginForm = document.getElementById("login");
    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;
    // console.log(formData);
    axios.post('/login', {
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

