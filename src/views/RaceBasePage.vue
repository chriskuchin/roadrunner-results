<template>
    <div>
        <breadcrumb :paths="paths"></breadcrumb>
        <router-view></router-view>
    </div>
</template>

<script>
import Breadcrumb from "../components/Breadcrumb.vue";
import { useRaceStore } from "../store/races";

const store = useRaceStore()
export default {
    components: {
        breadcrumb: Breadcrumb,
    },
    data: function () {
        return {
            race: {
                name: "",
                owner: "",
                id: "",
            }
        }
    },
    mounted: function () {
        store.loadRace(this.$route.params.raceId)
    },
    methods: {},
    computed: {
        paths: function () {
            var breadcrumb = [{
                path: "/races",
                display: "Races"
            }]
            var path = "/races"
            var pathSegments = this.$route.path.replace("/races/", "").split("/")
            for (let i = 0; i < pathSegments.length; i++) {
                var display = pathSegments[i]
                if (i == 0 && store.getName != "") {
                    display = store.getName
                }
                path += "/" + pathSegments[i]
                breadcrumb.push({
                    path: path,
                    display: display
                })
            }

            return breadcrumb;
        },
    },
};
</script>