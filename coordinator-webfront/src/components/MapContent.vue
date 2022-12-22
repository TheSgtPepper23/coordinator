<template>
  <v-container>
    <v-tabs v-model="tab" bg-color="white">
      <v-tab value="coords">Coordinates</v-tab>
      <v-tab value="milestones" :disabled="true">Milestones</v-tab>
    </v-tabs>

    <v-window v-model="tab">
      <v-window-item key="coord-wind" value="coords">
        <div class="coords-options">
          <v-btn style="margin-top: 30px;">Create new coordinate</v-btn>
          <v-text-field label="Filter coordinates" variant="solo" append-inner-icon="mdi-magnify"></v-text-field>
        </div>
        <div class="coord-container overflow-y-auto">
          <CoordinateItem v-for="coord in coordinates"></CoordinateItem>
        </div>
      </v-window-item>
    </v-window>
  </v-container>
</template>

<script>

import CoordinateItem from "./CoordinateItem.vue"

export default {
  data() {
    return {
      tab: null,
      coordinates: [1, 2, 3, 4, 54, 65,]
    }
  },
  components: { CoordinateItem }
}
</script>

<style scoped>
.coord-container {
  height: 79vh;
  margin-top: 30px;
}

.coords-options {
  display: flex !important;
  flex-direction: row !important;
}
</style>