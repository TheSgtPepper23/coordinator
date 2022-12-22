<template>
    <v-card :title="mapName" class="map-card" :class="{ 'map-selected': selectedMap == mapID }" @click="changeSelectedMap">
    </v-card>
</template>

<script>
import { useSelectedMapStore } from "@/stores/selectedMap"
import { mapWritableState } from "pinia";

export default {
    props: ["mapName", "mapID"],
    data() {
        return {
        };
    },
    methods: {
        changeSelectedMap() {
            this.selectedMap = this.mapID;
        }
    },
    computed: {
        ...mapWritableState(useSelectedMapStore, ["selectedMap"])
    },
    created() {
        this.mapStore = useSelectedMapStore();
    },

}
</script>

<style scoped>
.map-card {
    cursor: pointer;
    margin-bottom: 20px;
}

.map-selected {
    background-color: lightblue;
}

.map-card:hover {
    background-color: lightgray;
}
</style>