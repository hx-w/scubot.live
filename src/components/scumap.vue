<template>
  <div class="map-container">
    <amap
      cache-key="coord-picker-map"
      mmap-style="amap://styles/whitesmoke"
      async
      :center.sync="center"
      :zoom.sync="zoom"
      is-hotspot
      @click="onMapClick"
    >
      <amap-scale />
      <amap-tool-bar
        :position="{
          top: '110px',
          left: '10px',
        }"
      />
      <!-- <amap-marker v-if="position" :position.sync="position" draggable /> -->
      <amap-circle
        :center.sync="position"
        :radius.sync="radius"
        :fill-color="fill"
        :fill-opacity="0.5"
        :stroke-color="stroke"
        :editable="false"
        :draggable="true"
        stroke-style="dashed"
        :stroke-dasharray="[10, 10]"
      />
      <a-card
        :body-style="{
          'max-height': '450px',
          'overflow-y': 'scroll',
          'padding-top': 0,
        }"
        class="result-panel"
      >
        <template slot="title">
          <template v-if="mode === 'search'">
            <a-input-group compact style="display: flex">
              <a-auto-complete
                v-model="query"
                :data-source="tips"
                placeholder="输入关键词"
                style="flex: 1"
                @search="autoComplete"
              />
              <a-button @click="search(true)" :disabled="!query" type="primary">
                搜索
              </a-button>
            </a-input-group>
          </template>
          <template v-if="mode === 'result'">
            <div class="search-bar">
              <a-button icon="left" @click="reset" style="margin-right: 6px" />
              <span class="text"
                >搜索 {{ query }} 共
                {{ searching ? "..." : total }} 条结果</span
              >
            </div>
          </template>
        </template>

        <a-list
          v-if="mode === 'result'"
          :data-source="results"
          :loading="searching"
          item-layout="vertical"
          size="small"
          class="result-list"
        >
          <a-pagination
            slot="header"
            v-if="total > 0"
            v-model="pageIndex"
            :page-size="pageSize"
            :total="total"
            size="small"
          />
          <a-list-item slot="renderItem" slot-scope="item" :key="item.id">
            <a-list-item-meta :description="item.address">
              <span slot="title" style="cursor: pointer" @click="focus(item)">{{
                item.name
              }}</span>
            </a-list-item-meta>
          </a-list-item>
          <a-pagination
            slot="footer"
            v-if="total > 0"
            v-model="pageIndex"
            :page-size="pageSize"
            :total="total"
            size="small"
          />
        </a-list>
      </a-card>
    </amap>
  </div>
</template>

<script>
import { loadAmap, loadPlugins } from "@amap/amap-vue";

export default {
  name: 'scumap',
  props: {'icenter': Array, 'iposition': Array},
  data() {
    return {
      mode: "search",
      center: this.icenter,
      zoom: 16,
      query: "",
      searching: false,
      tips: [],
      results: [],
      position: this.iposition,
      pageIndex: 1,
      pageSize: 10,
      total: 0,
      radius: 50,
      fill: "#409EFF",
      stroke: "#000A58",
    };
  },
  computed: {
    wrapper() {
      return this.$refs.wrapper;
    },
    positionText() {
      if (!this.position) return "";
      return `${this.position[0]}, ${this.position[1]}`;
    },
  },
  created() {
    this.$ready = new Promise(async (resolve) => {
      const AMap = await loadAmap();
      await loadPlugins(["AMap.PlaceSearch", "AMap.AutoComplete"]);

      this.ps = new AMap.PlaceSearch({
        pageSize: this.pageSize,
      });
      this.ac = new AMap.AutoComplete();

      resolve();
    });
  },
  methods: {
    onMapClick(e) {
      if (e.lnglat) {
        this.position = [e.lnglat.lng, e.lnglat.lat];
      } else {
        this.position = null;
      }
    },
    async search(clear = false) {
      this.mode = "result";
      await this._ready;

      if (clear) {
        this.results = [];
        this.total = 0;
        this.pageIndex = 1;
        this.ps.setPageIndex(1);
      }

      this.searching = true;
      const { query } = this;
      this.ps.search(query, (status, result) => {
        this.searching = false;
        if (query !== this.query) return;

        if (status === "complete" && result.poiList) {
          this.results = result.poiList.pois;
          this.total = result.poiList.count;
        } else {
          this.results = [];
          this.total = 0;
        }
      });
    },
    async autoComplete(kw) {
      if (!kw) {
        this.tips = [];
      } else {
        this.ac.search(kw, (status, result) => {
          if (kw !== this.query) return;
          if (status === "complete" && result.tips) {
            this.tips = Array.from(new Set(result.tips.map((tip) => tip.name)));
          } else {
            this.tips = [];
          }
        });
      }
    },
    focus(poi) {
      const pos = [poi.location.lng, poi.location.lat];
      this.position = [...pos];
      this.center = [...pos];
    },
    reset() {
      this.ps.setPageIndex(1);
      this.results = [];
      this.tips = [];
      this.total = 0;
      this.mode = "search";
    },
  },
  watch: {
    pageIndex(value) {
      this.$ready.then(() => {
        this.ps.setPageIndex(value);
        this.search(false);
      });
    },
  },
};
</script>

<style scoped>
.map-container {
  width: 100%;
  height: 600px;
}

.result-panel {
  position: absolute;
  left: 10px;
  top: 10px;
  width: 320px;
  display: flex;
  flex-direction: column;

  .search-bar {
    display: flex;
    align-items: center;
    .text {
      text-overflow: ellipsis;
      flex: 1;
      overflow: hidden;
      white-space: normal;
    }
  }

  .result-list.ant-list-loading {
    min-height: 100px;
    display: flex;
    justify-content: center;
    align-items: center;
  }
}
.info {
  position: absolute;
  right: 10px;
  top: 60px;
  padding-left: 24px;
}
</style>