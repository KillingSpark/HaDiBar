/// <reference path="../node_modules/@types/jquery/index.d.ts" />

Vue.component('login-form', {
    data: function () {
        return {
            name: "",
            password: ""
        }
    },
    props: ['sessionid'],
    template: `
        <div class="navbar navbar-fixed-bottom">
            <span id="logintext">Log your ass in!</span>
            <input type=text v-model="name" placeholder="Name"/>
            <input type=text v-model="password" placeholder="Passwort"/>
            <button class="login_button" v-on:click="send_login">LOGIN</button>
        </div>
    `,
    methods: {
        send_login: function () {
            var comp = this
            $.ajax({
                url: "/login",
                data: { name: comp.name, password: comp.password },
                beforeSend: function (xhr) {
                    xhr.setRequestHeader("sessionID", comp.sessionid)
                },
                success: function(response){
                    alert(response)
                }
            })
        }
    }
})