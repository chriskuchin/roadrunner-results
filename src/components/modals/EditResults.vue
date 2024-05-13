<template>
  <modal ref="modal">
    <h1 class="title mb-4">Edit #{{ rowId + 1 }}</h1>
    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">From</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control">
            <input class="input" type="input" placeholder="Time" v-model="formattedResult">
          </p>
          <p class="help">{{ result }}</p>
        </div>
        <div class="field">
          <p class="control is-expanded">
            <input class="input" type="email" placeholder="Bib Number" v-model="bib">
          </p>
        </div>
      </div>
    </div>
    <div class="field is-grouped mt-6">
      <p class="control">
        <button class="button is-link" @click="updateResult">Save changes</button>
      </p>
      <p class="control">
        <button class="button" @click="close">Cancel</button>
      </p>
      <p class="control">
        <button class="button is-danger" @click="deleteResult">Delete Result</button>
      </p>
    </div>
  </modal>
</template>

<script>
import { mapActions } from 'pinia'
import { formatMilliseconds, calculateMilliseconds } from '../../utilities'
import { useResultsStore } from '../../store/results'
import Modal from '../Modal.vue'

export default {
  emits: ['reload'],
  components: {
    modal: Modal,
  },
  watch: {
    formattedResult: function (val) {
      this.result = calculateMilliseconds(val)
    }
  },
  data: function () {
    return {
      rowId: 0,
      bib: "",
      result: 0,
      formattedResult: ""
    }
  },
  methods: {
    ...mapActions(useResultsStore, ['deleteResultByRowId', 'getResultByRowId', 'updateResultByRowId']),
    open: function (rowId) {
      this.rowId = rowId

      const result = this.getResultByRowId(rowId)

      this.bib = result.bib_number
      this.result = result.result_ms
      this.formattedResult = formatMilliseconds(result.result_ms)

      this.$refs.modal.toggle()
    },
    close: function () {
      this.$refs.modal.toggle()
    },
    updateResult: function () {
      var that = this
      this.updateResultByRowId(this.rowId, this.result, this.bib).then(() => that.$emit('reload'))
      this.close()
    },
    deleteResult: function () {
      if (window.confirm("Are you sure you want to delete this result?")) {
        var that = this
        this.deleteResultByRowId(this.rowId).then(() => that.$emit('reload'))
        this.close()
      }
    },
  },
  computed: {
    formatResults() {
      return formatMilliseconds(this.result)
    },
  }
}
</script>