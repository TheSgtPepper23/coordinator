import { defineStore } from "pinia";

export const useSelectedMapStore = defineStore("selectedMap", {
  state: () => ({
    selectedMap: 0,
  }),
  actions: {
    changeSelectedMap(newMap) {
      this.selectedMap = newMap
    },
  },
});
