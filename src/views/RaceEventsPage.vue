<template>
    <div>
        <nav class="level">
            <div class="level-left">
                <div class="level-item">
                    <div class="field is-horizontal">
                        <div class="field-label is-small">
                            <label class="label">Birth Year</label>
                        </div>
                        <div class="select is-small">
                            <select v-model="filters.year">
                                <option value="0">None</option>
                                <option v-for="filter in birthYears" :key="filter" :value="filter">{{ filter }}
                                </option>
                            </select>
                        </div>
                    </div>
                </div>
                <div class="level-item">
                    <div class="field is-horizontal">
                        <div class="field-label is-small">
                            <label class="label">Gender</label>
                        </div>
                        <div class="select is-small">
                            <select v-model="filters.gender">
                                <option value="">None</option>
                                <option v-for="filter in genders" :key="filter" :value="filter">{{ filter }}
                                </option>
                            </select>
                        </div>
                    </div>
                </div>
            </div>
        </nav>
        <table class="table">
            <tr v-for="result in filtered" :key="result.bib_number">
                <td>{{ result.bib_number }}</td>
                <td>{{ result.first_name }}</td>
                <td>{{ result.last_name }}</td>
                <td>{{ result.birth_year }}</td>
                <td>{{ result.gender }}</td>
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
    data() {
        return {
            results: [],
            sortField: "result_ms",
            filters: {
                year: 0,
                gender: "",
            }
        }
    },
    computed: {
        birthYears: function () {
            return [...new Set(this.results.map(v => v.birth_year))].sort()
        },
        genders: function () {
            return [...new Set(this.results.map(v => v.gender))].sort()
        },
        filtered: function () {
            var results = []
            this.results.forEach(result => {
                if (this.unfiltered()) {
                    results.push(result)
                } else if (this.isfiltered(result)) {
                    results.push(result)
                }
            })
            return results.sort((a, b) => a[this.sortField] - b[this.sortField])
        },
    },
    methods: {
        unfiltered: function () {
            return this.filters.year == 0 && this.filters.gender == ""
        },
        isfiltered: function (result) {
            if (this.filters.year != 0 && this.filters.gender != "") {
                console.log("bother are set")
                return this.filters.year == result.birth_year && this.filters.gender == result.gender
            } else if (this.filters.year != 0) {
                console.log("year only")
                return this.filters.year == result.birth_year
            } else if (this.filters.gender != "") {
                console.log("only gender")
                return this.filters.gender == result.gender
            }

            return false
        },
        formatTime: function (ms) {
            return formatMilliseconds(ms)
        },
        async getResults(raceID, eventID) {
            let url = "/api/v1/races/" + raceID + "/events/" + eventID + "/results"
            this.results = await (await fetch(url)).json()
        }
    }

};
</script>