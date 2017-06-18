/// <reference path="node_modules/@types/jquery/index.d.ts" />

var bevapp = new Vue({
  el: '#bev-table',
  data: {
    sessionid: "Huhu",
    current_account: {
      Owner: {
        Name: "---",
      }, Value: 0
    },
    show_table: true,
    show_payment: false,
    beverages: [],
    accounts: []
  },
  computed: {
    bev_table: function (e) {
      return
    }
  },
  methods: {
    acc_selected: function (account) { this.current_account = account },
    make_payment: function (event) {
      var sum = 0
      for (var i = 0; i < this.beverages.length; i++) {
        sum += this.beverages[i].times * this.beverages[i].Value
      }
      this.changeAccount(-sum)
    },
    openApp: function (event, app_name) {
      var tabs = document.getElementsByClassName('tablink')
      for (var i = 0; i < tabs.length; i++) {
        tabs[i].classList.remove('active')
      }
      event.currentTarget.classList.add('active')
      if (app_name === 'bev-table') {
        this.show_table = true
        this.show_payment = false
      }
      if (app_name === 'direct-payment') {
        this.show_table = false
        this.show_payment = true
      }
    },
    changeAccount: function (diff) {
      var app = this
      $.ajax({
        url: "/account/update/" + app.current_account.ID,
        type: 'POST',
        data: { value: diff },
        beforeSend: function (xhr) {
          xhr.setRequestHeader("sessionID", app.sessionid)
        },
        success: function (response) {
          res = JSON.parse(response)
          if (res.status == "OK") {
            app.current_account.Value = res.response.Value
          }
        }
      })
    },
    updateAccounts: function () {
      var app = this
      $.ajax({
        url: "/account/all",
        type: 'GET',
        data: {},
        beforeSend: function (xhr) {
          xhr.setRequestHeader("sessionID", app.sessionid)
        },
        success: function (response) {
          res = JSON.parse(response)
          if (res.status == "OK") {
            app.accounts = res.response
          }
          app.current_account = app.accounts[0]
        }
      })
    },
    updateBeverages: function () {
      var app = this
      $.ajax({
        url: "/beverage/all",
        type: 'GET',
        data: {},
        beforeSend: function (xhr) {
          xhr.setRequestHeader("sessionID", app.sessionid)
        },
        success: function (response) {
          res = JSON.parse(response)
          if (res.status == "OK") {
            app.beverages = res.response
            app.beverages.forEach(function(element) {
              element.times = 0
            }, this);
          }
        }
      })
    },
  },
  created: function () {
    var app = this
    //getting an initial sessionID for the API
    $.get("/session", {}, function (data, textStatus, response) {
      app.sessionid = response.getResponseHeader('sessionID')
      app.updateAccounts()
      app.updateBeverages()
    })
  }
})



