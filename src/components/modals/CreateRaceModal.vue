<template>
  <modal ref="create-race-modal">
    <div class="tabs">
      <ul>
        <li :class="{ 'is-active': tabs.activeTab == 'create' }"><a @click="raceModalTabClick('create')">Create Race</a>
        </li>
        <li :class="{ 'is-active': tabs.activeTab == 'import' }"><a @click="raceModalTabClick('import')">Import Race</a>
        </li>
      </ul>
    </div>
    <div v-if="tabs.activeTab == 'create'">
      <p class="title mb-4 mt-2">Create Race</p>
      <div class="field mb-2">
        <label class="label">Description</label>
        <div class="control">
          <input class="input" type="text" placeholder="Race Description" v-model="description">
        </div>
      </div>
      <div class="field">
        <label class="label">Date</label>
        <div class="control">
          <input class="input" type="date" placeholder="" v-model="date">
        </div>
      </div>
      <div class="field is-grouped mt-4">
        <div class="control">
          <button :class="['button', 'is-link', { 'is-loading': creating }]" @click="createRace">Submit</button>
        </div>
        <div class="control">
          <button class="button is-link is-light" @click="cancelCreateRace">Cancel</button>
        </div>
      </div>
    </div>
    <div v-if="tabs.activeTab == 'import'">
      <p class="title mb-4 mt-2">Import Race</p>
      <div class="field mb-2">
        <label class="label">Race URL</label>
        <div class="control">
          <input class="input" type="text" placeholder="Race URL" v-model="importURL">
        </div>
      </div>
      <div class="field mb-2">
        <label class="label">Description</label>
        <div class="control">
          <input class="input" type="text" placeholder="Race Description" v-model="description">
        </div>
      </div>
      <div class="field">
        <label class="label">Date</label>
        <div class="control">
          <input class="input" type="date" placeholder="" v-model="date">
        </div>
      </div>
      <div class="field is-grouped mt-4">
        <div class="control">
          <button :class="['button', 'is-link', { 'is-loading': importing }]"
            @click="raceModalImportRace">Submit</button>
        </div>
        <div class="control">
          <button class="button is-link is-light" @click="cancelCreateRace">Cancel</button>
        </div>
      </div>
    </div>
  </modal>
</template>

<script>
import { mapStores } from 'pinia'
import Modal from '../Modal.vue'
import { createRace, importRace } from "../../api/races"
import { useRacesStore } from "../../store/races"

export default {
  components: {
    modal: Modal,
  },
  data: function () {
    return {
      creating: false,
      importing: false,
      description: "",
      date: "",
      importURL: "",
      tabs: {
        activeTab: "create",
      }
    }
  },
  methods: {
    open: function () {
      this.$refs['create-race-modal'].toggle()
    },
    close: function () {
      this.$refs['create-race-modal'].toggle()
    },
    raceModalImportRace: async function () {
      this.importing = true
      await importRace(this.importURL, this.description, this.date)
      this.close()
      this.racesStore.loadRaces()
      this.importing = false
    },
    raceModalTabClick: function (tabID) {
      this.tabs.activeTab = tabID
    },
    createRace: function () {
      var self = this
      self.creating = true
      createRace(self.description, self.date).then(() => {
        self.creating = false
        self.description = ""
        self.close()
        self.racesStore.loadRaces()
      })
    },
    cancelCreateRace: function () {
      this.$refs['create-race-modal'].toggle()
      this.description = ""
    },
  },
  computed: {
    ...mapStores(useRacesStore),
  }
}
</script>