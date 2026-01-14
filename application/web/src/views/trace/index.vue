<template>
  <div class="trace-container">
    <div class="trace-header">
      <div class="search-section">
        <el-input
          v-model="input"
          placeholder="请输入溯源码查询"
          style="width:300px;margin-right:15px;"
          clearable
        />
        <el-button type="primary" plain @click="FruitInfo">查询</el-button>
        <el-button
          type="success"
          plain
          icon="el-icon-camera"
          @click="showQRScanner = true"
        >
          二维码扫描
        </el-button>
        <el-button
          type="info"
          plain
          icon="el-icon-map-location"
          @click="showSupplyChainMap = true"
        >
          供应链地图
        </el-button>
      </div>
      <el-button type="success" plain @click="AllFruitInfo">获取所有农产品信息</el-button>
    </div>

    <!-- QR扫描 Dialog -->
    <el-dialog
      title="二维码扫描溯源"
      :visible.sync="showQRScanner"
      width="600px"
      @close="handleQRDialogClose"
    >
      <qr-scanner @scan-success="handleQRScanSuccess" />
    </el-dialog>

    <!-- 供应链地图 Dialog -->
    <el-dialog
      title="供应链可视化地图"
      :visible.sync="showSupplyChainMap"
      width="900px"
      top="5vh"
    >
      <supply-chain-map :trace-data="tracedata[0] || {}" />
    </el-dialog>

    <!-- 产品二维码 Dialog -->
    <product-qr-display
      :visible.sync="showProductQRDialog"
      :traceability-code="currentProduct?.traceability_code || ''"
      :product-name="currentProduct?.farmer_input?.fa_fruitName || ''"
      :product-info="currentProduct?.farmer_input || {}"
    />

    <el-table :data="tracedata" style="width:100%">
      <el-table-column label="溯源码">
        <template slot-scope="scope">
          <div class="trace-code-cell">
            <span>{{ scope.row.traceability_code }}</span>
            <el-button
              type="text"
              size="mini"
              class="qr-btn"
              @click="openProductQR(scope.row)"
            >
              <svg-icon icon-class="qrcode" /> 二维码
            </el-button>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="农产品名称" prop="farmer_input.fa_fruitName" />
      <el-table-column label="采摘时间" prop="farmer_input.fa_pickingTime" />

      <el-table-column label="操作" width="260" fixed="right">
        <template slot-scope="scope">
          <el-button
            size="mini"
            type="warning"
            icon="el-icon-warning"
            @click="handleDetectAnomalies(scope.row)"
          >异常检测</el-button>
          <el-button
            size="mini"
            type="danger"
            icon="el-icon-message-solid"
            @click="handleViewAlerts(scope.row)"
          >查看报警</el-button>
          <el-button
            size="mini"
            type="info"
            icon="el-icon-document"
            @click="handleViewRecalls(scope.row)"
          >召回记录</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 报警记录 Dialog -->
    <el-dialog title="报警记录" :visible.sync="showAlertsDialog" width="800px">
      <div v-loading="alertsLoading">
        <el-timeline v-if="currentAlerts.length">
          <el-timeline-item
            v-for="(alert,i) in currentAlerts"
            :key="i"
            :timestamp="alert.timestamp"
            :type="getSeverityType(alert.severity)"
          >
            <el-card>
              <p><strong>{{ alert.message }}</strong></p>
              <p>严重等级: {{ alert.severity }}/5</p>
            </el-card>
          </el-timeline-item>
        </el-timeline>
        <div v-else style="text-align:center;color:#999">暂无报警记录</div>
      </div>
    </el-dialog>

    <!-- 召回记录 Dialog -->
    <el-dialog title="召回记录" :visible.sync="showRecallsDialog" width="800px">
      <div v-loading="recallsLoading">
        <el-timeline v-if="currentRecalls.length">
          <el-timeline-item
            v-for="(recall,i) in currentRecalls"
            :key="i"
            type="danger"
            :timestamp="recall.timestamp"
          >
            <el-card>
              <p><strong>召回原因:</strong> {{ recall.reason }}</p>
            </el-card>
          </el-timeline-item>
        </el-timeline>
        <div v-else style="text-align:center;color:#999">暂无召回记录</div>
      </div>
    </el-dialog>

  </div>
</template>

<script>
import { getFruitInfo, getFruitList, getAllFruitInfo } from '@/api/trace'
import { detectAnomalies, getFruitAlerts, getFruitRecalls } from '@/api/alert'
import QRScanner from '@/components/QRScanner'
import SupplyChainMap from '@/components/SupplyChainMap'
import ProductQRDisplay from '@/components/ProductQRDisplay'

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
      currentRecalls: [],
      alertsLoading: false,
      recallsLoading: false
    }
  },
  created() {
    getFruitList().then(res => {
      this.tracedata = JSON.parse(res.data || '[]')
    })
  },
  methods: {
    FruitInfo() {
      if (!this.input.trim()) {
        this.$message.warning('请输入溯源码')
        return
      }
      const fd = new FormData()
      fd.append('traceability_code', this.input)
      getFruitInfo(fd).then(res => {
        this.tracedata = [JSON.parse(res.data)]
      })
    },

    AllFruitInfo() {
      getAllFruitInfo().then(res => {
        this.tracedata = JSON.parse(res.data || '[]')
      })
    },

    handleQRScanSuccess(code) {
      this.input = code
      this.showQRScanner = false
      this.FruitInfo()
    },

    handleQRDialogClose() {
      // 如果 QRScanner 组件内已处理 stopScanning，这里无需重复操作
      this.showQRScanner = false
    },

    openProductQR(product) {
      this.currentProduct = product
      this.showProductQRDialog = true
    },

    handleDetectAnomalies(row) {
      const fd = new FormData()
      fd.append('traceability_code', row.traceability_code)
      detectAnomalies(fd).then(() => {
        this.$message.success('异常检测完成')
      })
    },

    handleViewAlerts(row) {
      this.showAlertsDialog = true
      this.alertsLoading = true
      const fd = new FormData()
      fd.append('traceability_code', row.traceability_code)
      getFruitAlerts(fd).then(res => {
        this.currentAlerts = JSON.parse(res.data || '[]')
        this.alertsLoading = false
      })
    },

    handleViewRecalls(row) {
      this.showRecallsDialog = true
      this.recallsLoading = true
      const fd = new FormData()
      fd.append('traceability_code', row.traceability_code)
      getFruitRecalls(fd).then(res => {
        this.currentRecalls = JSON.parse(res.data || '[]')
        this.recallsLoading = false
      })
    },

    getSeverityType(level) {
      if (level >= 4) return 'danger'
      if (level >= 3) return 'warning'
      return 'success'
    }
  }
}
</script>

<style scoped>
.trace-container { margin:30px; }
.trace-code-cell { display:flex; justify-content:space-between; align-items:center; }
.qr-btn { font-size:12px; }
</style>
