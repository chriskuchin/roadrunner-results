<template>
    <div>
        <div class="box" v-for="event in raceStore.eventList" :key="event.id">
            <div class="title is-4">{{ event.description }}</div>
            <router-link :to="'/races/' + this.$route.params.raceId + '/events/' + event.eventId + '/results'"
                class="button">Event Results</router-link>
            {{ event }}
        </div>
    </div>
</template>

<script>
import { mapStores } from "pinia";
import { useRaceStore } from "../store/race";
export default {
    computed: {
        ...mapStores(useRaceStore)
    }
};
</script>