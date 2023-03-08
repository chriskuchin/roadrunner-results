<template>
    <div>
        {{ getName }}
        {{ yearValues }}
        {{ yearLabels }}
        <div style="width: 400px;">
            <canvas id="birth-year"></canvas>
        </div>
    </div>
</template>

<script>
import Chart from 'chart.js/auto'
import { useRaceStore } from "../store/races";
import { mapState } from "pinia";

const raceState = useRaceStore()

export default {
    mounted: function () {
        raceState.$subscribe(this.updateChart)
    },
    data: function () {
        return {
            chart: null
        }
    },
    methods: {
        updateChart: function () {
            if (this.chart == null) {
                this.chart = new Chart(document.getElementById('birth-year'),
                    {
                        type: 'bar',
                        data: {
                            labels: this.yearLabels,
                            datasets: [
                                {
                                    label: 'Athletes by birth year',
                                    data: this.yearValues
                                }
                            ]
                        }
                    }
                );
            }
        }
    },
    computed: {
        ...mapState(useRaceStore, ['getName', 'getID', 'yearValues', 'yearLabels'])
    }
};
</script>