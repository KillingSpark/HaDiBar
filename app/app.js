Vue.component('acc-option', {
  props: ['account', 'selected'],
  template: `<option v-on:click="selected(account)"> {{account.name}} </option>`
})

Vue.component('acc-select', {
  props: ['accs', 'selected'],
  template: `
    <select>
      <acc-option v-bind:selected="selected" v-for="acc in accs" v-bind:key="acc" v-bind:account="acc" />
    </select>
    `,
})

Vue.component('bev-table', {
  data: function () {
    return {
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
                </tr>
            </thead>
            <tbody>
                <tr v-for="(bev, index) in beverages" v-bind:bev="bev">
                    <td>{{bev.name}}</td>
                    <td>{{bev.value}}</td>
                    <td><input v-model="beverages[index].times" type="text" style="width: 100%" /></td>
                </tr>
            </tbody>
        </table>
      </div>
      <div class="row">
        <button v-on:click="exec">Execute</button>
      </div>
    </div>
    `
})

Vue.component('acc-info-table', {
  data: function () {
    return {
      difference: 0
    }
  },
  props: ['account', 'show_payment'],
  template: `
    <div>
      <table id="acc_table" class="table-bordered col-md-12">
          <thead>
              <tr>
                  <th>Name</th>
                  <th>Value</th>
              </tr>
          </thead>
          <tbody>
              <tr>
                  <td>{{account.name}}</td>
                  <td>{{account.value}}</td>
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
      this.account.value = Number(this.account.value) + Number(this.difference)
    }
  }
})

new Vue({
  el: '#bev-table',
  data: {
    current_account: { name: "--" },
    show_table: true,
    show_payment: false,
    beverages: [
      {
        name: "Bier",
        value: 1,
        times: 0
      },
      {
        name: "Cola",
        value: 2,
        times: 0
      }
    ],
    accounts: [
      {
        name: "Moritz",
        value: 1337
      },
      {
        name: "Paul",
        value: 1337
      },
      {
        name: "Chris",
        value: 1337
      }
    ]
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
        this.current_account.value -= this.beverages[i].times * this.beverages[i].value
      }
    },
    openApp: function (event, app_name) {
      var tabs = document.getElementsByClassName('tablink')
      for(var i = 0; i < tabs.length; i++){
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
    }
  }
})
