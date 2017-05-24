/// <reference path="../node_modules/@types/jquery/index.d.ts" />

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