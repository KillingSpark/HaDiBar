/// <reference path="node_modules/@types/jquery/index.d.ts" />

var bevapp = new Vue({
  el: '#bev-table',
  data: {
    token: "",
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
      $.post("/account/" + app.current_account.ID, { value: diff }, function (response) {
        app.current_account.Value = JSON.parse(response).Value
      })
    },
    updateAccounts: function () {
      var app = this
      $.get("/accounts", {}, function (response) {
        app.accounts = JSON.parse(response)
        app.current_account = app.accounts[0]
      })
    },
    updateBeverages: function () {
      var app = this
      $.get("/beverages", {}, function (response) {
        app.beverages = JSON.parse(response)
      })
    },
  },
  created: function () {
  }
})



