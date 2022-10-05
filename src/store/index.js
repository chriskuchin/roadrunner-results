import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)


const store = new Vuex.Store({
  state: {
  },
  mutations: {
    toggleMultiPrintMode(state) {
      state.print.multi = !state.print.multi
    },
  },
  getters: {
    multi_values: (state) => {
      return state.print.labelSettings.valueList
    },
  },
  actions: {}
})

export default store