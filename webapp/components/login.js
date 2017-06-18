/// <reference path="../node_modules/@types/jquery/index.d.ts" />

Vue.component('login-form', {
    data: function () {
        return {
            name: "",
            password: "",
            show_login: true
        }
    },
    props: ['sessionid'],
    template: `
        <div class="navbar navbar-fixed-bottom">
            <div v-if="show_login">
                <span id="logintext">Log your ass in!</span>
                <input type=text v-model="name" placeholder="Name"/>
                <input type=text v-model="password" placeholder="Passwort"/>
                <button class="login_button" v-on:click="send_login">LOGIN</button>
            </div>
            <div v-if="!show_login">
                <span id="logintext">Log your ass out!</span>
                <button class="login_button" v-on:click="send_logout">LOGOUT</button>
            </div>
        </div>
    `,
    methods: {
        send_login: function () {
            var comp = this
            $.ajax({
                url: "/session/login",
                type: "POST",
                data: { name: comp.name, password: comp.password },
                beforeSend: function (xhr) {
                    xhr.setRequestHeader("sessionID", comp.sessionid)
                },
                success: function(response){
                    if(response === "OK"){
                        comp.show_login = false
                    }else{
                        alert("Login failed")
                    }
                }
            })
        },
        send_logout: function () {
            var comp = this
            $.ajax({
                url: "/session/logout",
                type: "POST",
                data: {},
                beforeSend: function (xhr) {
                    xhr.setRequestHeader("sessionID", comp.sessionid)
                },
                success: function(response){
                    if(response === "OK"){
                        comp.show_login = true
                    }else{
                        alert("Logout failed")
                    }
                }
            })
        }
    }
})