<template>
  <v-container :fluid="true" class="bordered-sidebar">
    <vue-final-modal v-model="creationModal" :click-to-close="false" :drag="false">
      <v-card class="modal-card" title="New map">
        <v-form ref="creationForm" v-model="creationValid" class="creation-form">
          <v-text-field label="Map name" variant="solo" required v-model="creationName"
            :rules="[(v) => !!v || 'Name is required']"></v-text-field>
          <v-select label="Map version" :items="['Java', 'Bedrock']" variant="solo" required v-model="creationVersion"
            :rules="[(v) => !!v || 'Version is required']"></v-select>
        </v-form>
        <div class="button-div">
          <v-btn @click="cancelCreation">Cancel</v-btn>
          <v-btn @click="createNewMap">Save</v-btn>
        </div>
      </v-card>
    </vue-final-modal>
    <v-text-field label="Filtrar mapa" variant="solo" append-inner-icon="mdi-magnify"></v-text-field>
    <div class="map-container overflow-y-auto">
      <MapItem v-for="mapObj in maps" :mapName="mapObj.name" :mapID="mapObj.id"></MapItem>
    </div>
    <div class="button-div">
      <v-btn icon="mdi-plus" color="white" @click="creationModal = true"></v-btn>
      <v-btn icon="mdi-pencil" color="white"></v-btn>
      <v-btn icon="mdi-delete" color="white" @click="deleteSelectedMap"></v-btn>
    </div>
  </v-container>
</template>
<script>
import axios from "axios";
import MapItem from "./MapItem.vue";
import { useSelectedMapStore } from "@/stores/selectedMap";
import { mapWritableState } from "pinia";

export default {
  data() {
    return {
      maps: [],
      url: "",
      creationModal: false,
      creationValid: true,
      creationName: "",
      creationVersion: "",
    };
  },
  methods: {
    deleteSelectedMap() {
      axios.delete(`${this.url}map/${this.selectedMap}`).then((data) => {
        this.getMaps();
        console.log("The map has been deleted");
      });
    },
    getMaps() {
      axios
        .get(`${this.url}map/`)
        .then((data) => {
          this.maps = data.data.body;
          if (this.maps.length > 0) {
            this.selectedMap = this.maps[0].id;
          }
        })
        .catch((err) => {
          console.log("An error ocurred while getting the maps");
        });
    },
    async createNewMap() {
      const { valid } = await this.$refs.creationForm.validate();

      if (valid) {
        let body = {
          name: this.creationName,
          version: this.creationVersion,
        };

        axios.post(`${this.url}map/`, body).then((resp) => {
          this.$refs.creationForm.reset();

          this.creationModal = false;
          this.getMaps();
        });
      }
    },
    cancelCreation() {
      this.$refs.creationForm.reset();
      this.creationModal = false;
    },
  },
  computed: {
    ...mapWritableState(useSelectedMapStore, ["selectedMap"]),
  },
  created() {
    this.url = import.meta.env.VITE_URL;
  },
  mounted() {
    this.getMaps();
  },
  components: { MapItem },
};
</script>

<style scoped>
.bordered-sidebar {
  text-align: left;
}

.map-container {
  height: 80vh;
}

.button-div {
  display: flex;
  justify-content: space-evenly;
  margin-top: 10px;
}

.modal-card {
  width: 400px;
  height: 300px;
  margin: 0;
  position: absolute;
  top: 50%;
  left: 50%;
  -ms-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
}

.creation-form {
  padding: 20px;
}
</style>
