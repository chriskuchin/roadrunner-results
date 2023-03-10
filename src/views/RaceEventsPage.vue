<template>
    <div>
        <table class="table">
            <tr v-for="result in results" :key="result.bib_number">
                <td>{{ result.bib_number }}</td>
                <td>{{ result.first_name }}</td>
                <td>{{ result.last_name }}</td>
                <td>{{ result.birth_year }}</td>
                <td>{{ result.gender }}</td>
                <td>{{ result.result_ms }}</td>
                <td>{{ formatTime(result.result_ms) }}</td>
            </tr>
        </table>
    </div>
</template>

<script>
import { formatMilliseconds } from "../utilities"

export default {
    mounted: function () {
        this.getResults(this.$route.params["raceId"], this.$route.params["eventId"])
    },
    data: function () {
        return {
            results: []
        }
    },
    methods: {
        formatTime: function (ms) {
            return formatMilliseconds(ms)
        },
        async getResults(raceID, eventID) {
            let url = "/api/v1/races/" + raceID + "/events/" + eventID + "/results"
            let results = await (await fetch(url)).json()

            results.sort((a, b) => a.result_ms - b.result_ms)

            this.results = results
        }
    }

};
</script>