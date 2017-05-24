/// <reference path="../node_modules/@types/jquery/index.d.ts" />

Vue.component('bev-table', {
  data: function () {
    return {
      bev_value: 0,
      bev_name: ""
    }
  },
  props: ['beverages', 'exec', 'sessionid'],
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
        <div class="col-md-6" id="makedrink">
            <span id="makedrinktext">MAKE ALL THE DRINKS!</span>
            <br>
            <input type="text" v-model="bev_name" placeholder="new name" />
            <input type="text" v-model="bev_value" placeholder="new value" />
            <br>
            <button style="margin-top: 1%;" class="col-md-2" v-on:click="addBeverage">Add</button>
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
        data: {},
        beforeSend: function (xhr) {
          xhr.setRequestHeader("sessionID", comp.sessionid)
        },
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
        beforeSend: function (xhr) {
          xhr.setRequestHeader("sessionID", comp.sessionid)
        },
        success: function (response) {
          comp.beverages.push(JSON.parse(response))
        }
      });
    }
  }
})
