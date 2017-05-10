Vue.component('bev-table', {
  data: function () {
    return {
      bev_value: 0,
      bev_name: ""
    }
  },
  props: ['beverages', 'exec'],
  template:
  ` <div>
      <div class="row">
        <table id="bev_table" class="table-bordered table-hover col-md-3">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Value</th>
                    <th>Times</th>
                    <th>Delete</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(bev, index) in beverages" v-bind:bev="bev">
                    <td>{{bev.Name}}</td>
                    <td class="center">{{bev.Value}}</td>
                    <td><input v-model="beverages[index].times" type="text" style="width: 100%" /></td>
                    <td class="danger center" v-on:click="deleteBeverage(index)">X</td>
                </tr>
            </tbody>
        </table>
        <div class="col-md-8" style="float:right;">
          <h4>MAKE ALL THE DRINKS!</h4>
          <br>
          <input class="col-md-3" type="text" v-model="bev_name" placeholder="new name" />
          <input class="col-md-3" type="text" v-model="bev_value" placeholder="new value" />
          <button class="col-md-2" v-on:click="addBeverage">Add</button>
        </div>
      </div>
      <div class="row">
        <button v-on:click="exec">Execute</button>
      </div>
    </div>
    `,
  methods: {
    deleteBeverage: function (index) {
      var comp = this
      $.ajax({
        url: "/beverage/" + comp.beverages[index].ID,
        type: 'DELETE',
        success: function (response) {
          comp.beverages.splice(index, 1)
        }
      });
    },
    addBeverage: function () {
      var comp = this
      $.ajax({
        url: "/newbeverage",
        type: 'PUT',
        data: { name: this.bev_name, value: this.bev_value },
        success: function (response) {
          comp.beverages.push(JSON.parse(response))
        }
      });
    }
  }
})

Vue.component('acc-option', {
  props: ['account', 'selected'],
  template: `<option v-on:click="selected(account)"> {{account.Owner.Name}} </option>`
})

Vue.component('acc-select', {
  props: ['accs', 'selected'],
  template: `
    <select style="width: 100%">
      <acc-option v-bind:selected="selected" v-for="acc in accs" v-bind:key="acc" v-bind:account="acc" />
    </select>
    `,
})

Vue.component('acc-info-table', {
  data: function () {
    return {
      difference: 0
    }
  },
  computed: {
    isNegativ: function () {
      return Number(this.account.Value) < 0
    }
  },
  props: ['account', 'show_payment', 'accs', 'selected'],
  template: `
    <div>
      <table id="acc_table" class="table-bordered col-md-3">
          <thead>
              <tr>
                  <th>Name</th>
                  <th>Value</th>
              </tr>
          </thead>
          <tbody>
              <tr>
                  <td><acc-select v-bind:accs="accs" v-bind:selected="selected"></acc-select></td>
                  <td v-bind:class="{danger: isNegativ, success: !isNegativ}">{{account.Value}}</td>
              </tr>
          </tbody>
      </table>
      <div v-if="show_payment">
        <input type="text" v-model="difference" />
        <button v-on:click="make_payment">Execute</button>
      </div>
    </div>
  `,
  methods: {
    make_payment: function () {
      this.account.Value = Number(this.account.Value) + Number(this.difference)
    }
  }
})

var bevapp = new Vue({
  el: '#bev-table',
  data: {
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
      $.post("/account/" + app.current_account.ID,{value: diff}, function(response){
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
    }
  },
  created: function () {
    this.updateBeverages()
    this.updateAccounts()
  }
})



