
function postLogin(username, password) {
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

