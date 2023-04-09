<template>
    <div>
        list heats here
        <div class="box" v-for="timer in eventStore.timerList">
            {{ timer }}
        </div>
    </div>
</template>

<script>
import { useEventStore } from "../store/event";
import { mapStores } from "pinia";
export default {
    components: {
    },
    data: function () {
        return {};
    },
    mounted: function () {
        console.log("mounted")
        this.eventStore.loadTimers(this.$route.params.raceId, this.$route.params.eventId)
    },
    actions: {

    },
    computed: {
        ...mapStores(useEventStore)
    },
};
</script>