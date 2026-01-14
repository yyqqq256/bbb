<template>
  <div class="trace-container">

    <div class="trace-header">
      <div class="search-section">
        <el-input v-model="input" placeholder="请输入溯源码查询" style="width:300px;margin-right:15px"/>
        <el-button type="primary" plain @click="FruitInfo">查询</el-button>
        <el-button type="success" plain @click="showQRScanner = true">二维码扫描</el-button>
        <el-button type="info" plain @click="showSupplyChainMap = true">供应链地图</el-button>
      </div>
      <el-button type="success" plain @click="AllFruitInfo">获取所有农产品信息</el-button>
    </div>

    <!-- 扫码 -->
    <el-dialog title="二维码扫描溯源" :visible.sync="showQRScanner" width="600px">
      <qr-scanner @scan-success="handleQRScanSuccess"/>
    </el-dialog>

    <!-- 地图 -->
    <el-dialog title="供应链可视化地图" :visible.sync="showSupplyChainMap" width="900px" top="5vh">
      <supply-chain-map :trace-data="tracedata.length ? tracedata[0] : {}"/>
    </el-dialog>

    <!-- 产品二维码 -->
    <product-qr-display
      :visible.sync="showProductQRDialog"
      :traceability-code="currentProduct ? currentProduct.traceability_code : ''"
      :product-name="currentProduct && currentProduct.farmer_input ? currentProduct.farmer_input.fa_fruitName : ''"
      :product-info="currentProduct && currentProduct.farmer_input ? currentProduct.farmer_input : {}"
    />

    <!-- 表格 -->
    <el-table :data="tracedata" style="width:100%">

      <el-table-column label="溯源码">
        <template slot-scope="scope">
          <div class="trace-code-cell">
            <span>{{ scope.row.traceability_code }}</span>
            <el-button type="text" size="mini" @click="openProductQR(scope.row)">二维码</el-button>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="农产品名称" prop="farmer_input.fa_fruitName"/>
      <el-table-column label="采摘时间" prop="farmer_input.fa_pickingTime"/>

      <el-table-column label="操作" width="260" fixed="right">
        <template slot-scope="scope">
          <el-button size="mini" type="warning" @click="handleDetectAnomalies(scope.row)">异常检测</el-button>
          <el-button size="mini" type="danger" @click="handleViewAlerts(scope.row)">报警</el-button>
          <el-button size="mini" type="info" @click="handleViewRecalls(scope.row)">召回</el-button>
        </template>
      </el-table-column>

    </el-table>

    <!-- 报警 -->
    <el-dialog title="报警记录" :visible.sync="showAlertsDialog" width="800px">
      <el-timeline v-if="currentAlerts.length">
        <el-timeline-item
          v-for="(alert,i) in currentAlerts"
          :key="i"
          :timestamp="alert.timestamp">
          {{ alert.message }}
        </el-timeline-item>
      </el-timeline>
      <div v-else style="text-align:center;color:#999">暂无报警</div>
    </el-dialog>

    <!-- 召回 -->
    <el-dialog title="召回记录" :visible.sync="showRecallsDialog" width="800px">
      <el-timeline v-if="currentRecalls.length">
        <el-timeline-item
          v-for="(recall,i) in currentRecalls"
          :key="i"
          :timestamp="recall.timestamp">
          {{ recall.reason }}
        </el-timeline-item>
      </el-timeline>
      <div v-else style="text-align:center;color:#999">暂无召回</div>
    </el-dialog>

  </div>
</template>

<script>
import QRScanner from '@/components/QRScanner'
import SupplyChainMap from '@/components/SupplyChainMap'
import ProductQRDisplay from '@/components/ProductQRDisplay'

import {
  getFruitInfo,
  getFruitList,
  getAllFruitInfo
} from '@/api/trace'

import {
  detectAnomalies,
  getFruitAlerts,
  getFruitRecalls
} from '@/api/alert'

export default {
  name: 'Trace',
  components: { QRScanner, SupplyChainMap, ProductQRDisplay },

  data() {
    return {
      tracedata: [],
      input: '',
      showQRScanner: false,
      showSupplyChainMap: false,
      showProductQRDialog: false,
      currentProduct: null,
      showAlertsDialog: false,
      showRecallsDialog: false,
      currentAlerts: [],
      currentRecalls: []
    }
  },

  created() {
    getFruitList().then(res => {
      this.tracedata = JSON.parse(res.data)
    })
  },

  methods: {
    AllFruitInfo() {
      getAllFruitInfo().then(res => {
        this.tracedata = JSON.parse(res.data)
      })
    },

    FruitInfo() {
      const fd = new FormData()
      fd.append('traceability_code', this.input)
      getFruitInfo(fd).then(res => {
        this.tracedata = [JSON.parse(res.data)]
      })
    },

    handleQRScanSuccess(code) {
      this.input = code
      this.showQRScanner = false
      this.FruitInfo()
    },

    openProductQR(row) {
      this.currentProduct = row
      this.showProductQRDialog = true
    },

    handleDetectAnomalies(row) {
      const fd = new FormData()
      fd.append('traceability_code', row.traceability_code)
      detectAnomalies(fd)
    },

    handleViewAlerts(row) {
      this.showAlertsDialog = true
      const fd = new FormData()
      fd.append('traceability_code', row.traceability_code)
      getFruitAlerts(fd).then(res => {
        this.currentAlerts = JSON.parse(res.data || '[]')
      })
    },

    handleViewRecalls(row) {
      this.showRecallsDialog = true
      const fd = new FormData()
      fd.append('traceability_code', row.traceability_code)
      getFruitRecalls(fd).then(res => {
        this.currentRecalls = JSON.parse(res.data || '[]')
      })
    }
  }
}
</script>

<style scoped>
.trace-container {
  margin: 30px;
}
.trace-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}
.search-section {
  display: flex;
  align-items: center;
}
.trace-code-cell {
  display: flex;
  justify-content: space-between;
}
</style>
