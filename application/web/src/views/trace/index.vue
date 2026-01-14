<template>
  <div class="trace-container">
    <div class="trace-header">
      <div class="search-section">
        <el-input v-model="input" placeholder="请输入溯源码查询" style="width: 300px;margin-right: 15px;" />
        <el-button type="primary" plain @click="FruitInfo"> 查询 </el-button>
        <el-button type="success" plain @click="showQRScanner = true" icon="el-icon-camera">
          二维码扫描
        </el-button>
        <el-button type="info" plain @click="showSupplyChainMap = true" icon="el-icon-map-location">
          供应链地图
        </el-button>
      </div>
      <el-button type="success" plain @click="AllFruitInfo"> 获取所有农产品信息 </el-button>
    </div>
    
    <!-- QR扫描对话框 -->
    <el-dialog
      title="二维码扫描溯源"
      :visible.sync="showQRScanner"
      width="600px"
      :before-close="handleQRScannerClose"
    >
      <qr-scanner @scan-success="handleQRScanSuccess" @close="showQRScanner = false" />
    </el-dialog>
    
    <!-- 供应链地图对话框 -->
    <el-dialog
      title="供应链可视化地图"
      :visible.sync="showSupplyChainMap"
      width="900px"
      top="5vh"
    >
      <supply-chain-map :trace-data="tracedata[0] || {}" />
    </el-dialog>
    
    <!-- 产品二维码显示对话框 -->
    <product-qr-display
      :visible.sync="showProductQR"
      :traceability-code="currentProduct ? currentProduct.traceability_code : ''"
      :product-name="currentProduct ? currentProduct.farmer_input.fa_fruitName : ''"
      :product-info="currentProduct ? currentProduct.farmer_input : {}"
    />
    <el-table
      :data="tracedata"
      style="width: 100%"
    >
      <el-table-column type="expand">
        <template slot-scope="props">
          <div class="trace-timeline-container">
            <el-timeline>
              <!-- 种植端 -->
              <el-timeline-item
                :timestamp="props.row.farmer_input.fa_timestamp"
                placement="top"
                color="#67C23A"
                icon="el-icon-s-custom"
                size="large"
              >
                <el-card shadow="hover" class="trace-card">
                  <div slot="header" class="clearfix">
                    <span class="card-title" style="color: #67C23A;">
                      <i class="el-icon-s-custom"></i> 种植端信息
                    </span>
                    <el-tag size="small" type="success" style="float: right;">交易ID: {{ props.row.farmer_input.fa_txid ? props.row.farmer_input.fa_txid.substring(0, 10) + '...' : '暂无' }}</el-tag>
                  </div>
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">农产品名称：</span>
                        <span class="value">{{ props.row.farmer_input.fa_fruitName }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">产地：</span>
                        <span class="value">{{ props.row.farmer_input.fa_origin }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">种植户名称：</span>
                        <span class="value">{{ props.row.farmer_input.fa_farmerName }}</span>
                      </div>
                    </el-col>
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">种植时间：</span>
                        <span class="value">{{ props.row.farmer_input.fa_plantTime }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">采摘时间：</span>
                        <span class="value">{{ props.row.farmer_input.fa_pickingTime }}</span>
                      </div>
                    </el-col>
                  </el-row>
                  <div v-if="props.row.farmer_input.fa_img" class="image-section">
                    <span class="label">相关图片：</span>
                    <a :href="`${baseApi}getImg/${props.row.farmer_input.fa_img}`" target="_blank">
                      <el-image
                        style="width: 100px; height: 100px; border-radius: 4px;"
                        :src="`${baseApi}getImg/${props.row.farmer_input.fa_img}`"
                        fit="cover"
                      />
                    </a>
                  </div>
                </el-card>
              </el-timeline-item>

              <!-- 加工端 -->
              <el-timeline-item
                :timestamp="props.row.factory_input.fac_timestamp"
                placement="top"
                color="#409EFF"
                icon="el-icon-s-cooperation"
                size="large"
              >
                <el-card shadow="hover" class="trace-card">
                  <div slot="header" class="clearfix">
                    <span class="card-title" style="color: #409EFF;">
                      <i class="el-icon-s-cooperation"></i> 工厂端信息
                    </span>
                    <el-tag size="small" style="float: right;">交易ID: {{ props.row.factory_input.fac_txid ? props.row.factory_input.fac_txid.substring(0, 10) + '...' : '暂无' }}</el-tag>
                  </div>
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">商品名称：</span>
                        <span class="value">{{ props.row.factory_input.fac_productName }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">生产批次：</span>
                        <span class="value">{{ props.row.factory_input.fac_productionbatch }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">工厂名称：</span>
                        <span class="value">{{ props.row.factory_input.fac_factoryName }}</span>
                      </div>
                    </el-col>
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">生产时间：</span>
                        <span class="value">{{ props.row.factory_input.fac_productionTime }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">联系电话：</span>
                        <span class="value">{{ props.row.factory_input.fac_contactNumber }}</span>
                      </div>
                    </el-col>
                  </el-row>
                  <div v-if="props.row.factory_input.fac_img" class="image-section">
                    <span class="label">相关图片：</span>
                    <a :href="`${baseApi}getImg/${props.row.factory_input.fac_img}`" target="_blank">
                      <el-image
                        style="width: 100px; height: 100px; border-radius: 4px;"
                        :src="`${baseApi}getImg/${props.row.factory_input.fac_img}`"
                        fit="cover"
                      />
                    </a>
                  </div>
                </el-card>
              </el-timeline-item>

              <!-- 物流端 -->
              <el-timeline-item
                :timestamp="props.row.driver_input.dr_timestamp"
                placement="top"
                color="#E6A23C"
                icon="el-icon-truck"
                size="large"
              >
                <el-card shadow="hover" class="trace-card">
                  <div slot="header" class="clearfix">
                    <span class="card-title" style="color: #E6A23C;">
                      <i class="el-icon-truck"></i> 物流端信息
                    </span>
                    <el-tag size="small" type="warning" style="float: right;">交易ID: {{ props.row.driver_input.dr_txid ? props.row.driver_input.dr_txid.substring(0, 10) + '...' : '暂无' }}</el-tag>
                  </div>
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">司机姓名：</span>
                        <span class="value">{{ props.row.driver_input.dr_name }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">车牌号：</span>
                        <span class="value">{{ props.row.driver_input.dr_carNumber }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">联系电话：</span>
                        <span class="value">{{ props.row.driver_input.dr_phone }}</span>
                      </div>
                    </el-col>
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">年龄：</span>
                        <span class="value">{{ props.row.driver_input.dr_age }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">运输记录：</span>
                        <span class="value">{{ props.row.driver_input.dr_transport }}</span>
                      </div>
                    </el-col>
                  </el-row>
                  <div v-if="props.row.driver_input.dr_img" class="image-section">
                    <span class="label">相关图片：</span>
                    <a :href="`${baseApi}getImg/${props.row.driver_input.dr_img}`" target="_blank">
                      <el-image
                        style="width: 100px; height: 100px; border-radius: 4px;"
                        :src="`${baseApi}getImg/${props.row.driver_input.dr_img}`"
                        fit="cover"
                      />
                    </a>
                  </div>
                </el-card>
              </el-timeline-item>

              <!-- 销售端 -->
              <el-timeline-item
                :timestamp="props.row.shop_input.sh_timestamp"
                placement="top"
                color="#909399"
                icon="el-icon-s-shop"
                size="large"
              >
                <el-card shadow="hover" class="trace-card">
                  <div slot="header" class="clearfix">
                    <span class="card-title" style="color: #909399;">
                      <i class="el-icon-s-shop"></i> 销售端信息
                    </span>
                    <el-tag size="small" type="info" style="float: right;">交易ID: {{ props.row.shop_input.sh_txid ? props.row.shop_input.sh_txid.substring(0, 10) + '...' : '暂无' }}</el-tag>
                  </div>
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">商店名称：</span>
                        <span class="value">{{ props.row.shop_input.sh_shopName }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">商店位置：</span>
                        <span class="value">{{ props.row.shop_input.sh_shopAddress }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">联系电话：</span>
                        <span class="value">{{ props.row.shop_input.sh_shopPhone }}</span>
                      </div>
                    </el-col>
                    <el-col :span="12">
                      <div class="detail-item">
                        <span class="label">入库时间：</span>
                        <span class="value">{{ props.row.shop_input.sh_storeTime }}</span>
                      </div>
                      <div class="detail-item">
                        <span class="label">销售时间：</span>
                        <span class="value">{{ props.row.shop_input.sh_sellTime }}</span>
                      </div>
                    </el-col>
                  </el-row>
                  <div v-if="props.row.shop_input.sh_img" class="image-section">
                    <span class="label">相关图片：</span>
                    <a :href="`${baseApi}getImg/${props.row.shop_input.sh_img}`" target="_blank">
                      <el-image
                        style="width: 100px; height: 100px; border-radius: 4px;"
                        :src="`${baseApi}getImg/${props.row.shop_input.sh_img}`"
                        fit="cover"
                      />
                    </a>
                  </div>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        label="溯源码"
        prop="traceability_code"
      >
        <template slot-scope="scope">
          <div class="trace-code-cell">
            <span>{{ scope.row.traceability_code }}</span>
            <el-button
              type="text"
              size="mini"
              @click="showProductQR(scope.row)"
              class="qr-btn"
            >
              <svg-icon icon-class="qrcode" style="margin-right: 2px;" />
              二维码
            </el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column
        label="操作"
        width="250"
        fixed="right"
      >
        <template slot-scope="scope">
          <el-button
            type="warning"
            size="mini"
            @click="handleDetectAnomalies(scope.row)"
            icon="el-icon-warning"
          >
            异常检测
          </el-button>
          <el-button
            type="danger"
            size="mini"
            @click="handleViewAlerts(scope.row)"
            icon="el-icon-message-solid"
          >
            查看报警
          </el-button>
          <el-button
            type="info"
            size="mini"
            @click="handleViewRecalls(scope.row)"
            icon="el-icon-document"
          >
            召回记录
          </el-button>
        </template>
      </el-table-column>
      <el-table-column
        label="农产品名称"
        prop="farmer_input.fa_fruitName"
      />
      <el-table-column
        label="农产品采摘时间"
        prop="farmer_input.fa_pickingTime"
      />
    </el-table>
    
    <!-- 报警记录对话框 -->
    <el-dialog
      title="报警记录"
      :visible.sync="showAlertsDialog"
      width="800px"
      top="5vh"
    >
      <div v-loading="alertsLoading">
        <div v-if="currentAlerts.length === 0" style="text-align: center; padding: 20px; color: #909399;">
          暂无报警记录
        </div>
        <el-timeline v-else>
          <el-timeline-item
            v-for="(alert, index) in currentAlerts"
            :key="index"
            :timestamp="alert.timestamp"
            :type="getSeverityType(alert.severity)"
            placement="top"
          >
            <el-card>
              <div slot="header" class="clearfix">
                <span style="font-weight: bold;">{{ alert.message }}</span>
                <el-tag
                  :type="getSeverityType(alert.severity)"
                  size="mini"
                  style="float: right;"
                >
                  严重等级: {{ alert.severity }}/5
                </el-tag>
              </div>
              <div class="alert-detail">
                <p><strong>报警ID:</strong> {{ alert.alert_id }}</p>
                <p><strong>异常类型:</strong> {{ alert.anomaly_type }}</p>
                <p><strong>状态:</strong>
                  <el-tag :type="getStatusType(alert.status)" size="mini">
                    {{ alert.status === 'active' ? '活跃' : alert.status === 'resolved' ? '已解决' : '待处理' }}
                  </el-tag>
                </p>
                <p v-if="alert.description"><strong>描述:</strong> {{ alert.description }}</p>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
    
    <!-- 召回记录对话框 -->
    <el-dialog
      title="召回记录"
      :visible.sync="showRecallsDialog"
      width="800px"
      top="5vh"
    >
      <div v-loading="recallsLoading">
        <div v-if="currentRecalls.length === 0" style="text-align: center; padding: 20px; color: #909399;">
          暂无召回记录
        </div>
        <el-timeline v-else>
          <el-timeline-item
            v-for="(recall, index) in currentRecalls"
            :key="index"
            :timestamp="recall.timestamp"
            type="danger"
            placement="top"
          >
            <el-card>
              <div slot="header" class="clearfix">
                <span style="font-weight: bold; color: #F56C6C;">产品召回</span>
                <el-tag
                  :type="recall.status === 'active' ? 'danger' : 'success'"
                  size="mini"
                  style="float: right;"
                >
                  {{ recall.status === 'active' ? '进行中' : '已完成' }}
                </el-tag>
              </div>
              <div class="recall-detail">
                <p><strong>召回ID:</strong> {{ recall.recall_id }}</p>
                <p><strong>召回原因:</strong> {{ recall.reason }}</p>
                <p><strong>相关报警ID:</strong> {{ recall.alert_id }}</p>
                <p v-if="recall.description"><strong>详细说明:</strong> {{ recall.description }}</p>
                <p><strong>批次范围:</strong> {{ recall.batch_range || '全部批次' }}</p>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { getFruitInfo, getFruitList, getAllFruitInfo, getFruitHistory } from '@/api/trace'
import { detectAnomalies, getFruitAlerts, getFruitRecalls } from '@/api/alert'
import QRScanner from '@/components/QRScanner'
import SupplyChainMap from '@/components/SupplyChainMap'
import ProductQRDisplay from '@/components/ProductQRDisplay'

export default {
  name: 'Trace',
  components: {
    QRScanner,
    SupplyChainMap,
    ProductQRDisplay
  },
  data() {
    return {
      tracedata: [],
      loading: false,
      input: '',
      baseApi: process.env.VUE_APP_BASE_API,
      showQRScanner: false,
      showSupplyChainMap: false,
      showProductQR: false,
      currentProduct: null,
      showAlertsDialog: false,
      showRecallsDialog: false,
      currentAlerts: [],
      currentRecalls: [],
      alertsLoading: false,
      recallsLoading: false
    }
  },
  computed: {
    ...mapGetters([
      'name',
      'userType'
    ])
  },
  created() {
    const code = this.$route.params.traceability_code
    if (code) {
      this.input = code
      this.FruitInfo()
    } else {
      getFruitList().then(res => {
        this.tracedata = JSON.parse(res.data)
      })
    }
  },
  methods: {
    AllFruitInfo() {
      getAllFruitInfo().then(res => {
        this.tracedata = JSON.parse(res.data)
      })
    },
    FruitHistory() {
      getFruitHistory().then(res => {
        // console.log(res)
      })
    },
    FruitInfo() {
      var formData = new FormData()
      formData.append('traceability_code', this.input)
      getFruitInfo(formData).then(res => {
        if (res.code === 200) {
          // console.log(res)
          this.tracedata = []
          this.tracedata[0] = JSON.parse(res.data)
          return
        } else {
          this.$message.error(res.message)
        }
      })
    },
    
    handleQRScanSuccess(qrCode) {
      this.input = qrCode
      this.showQRScanner = false
      this.FruitInfo()
      this.$message.success(`二维码扫描成功: ${qrCode}`)
    },
    
    handleQRScannerClose() {
      this.showQRScanner = false
    },
    
    showProductQR(product) {
      this.currentProduct = product
      this.showProductQR = true
    },
    
    // 异常检测
    handleDetectAnomalies(row) {
      this.$confirm('确定要对该产品进行异常检测吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const formData = new FormData()
        formData.append('traceability_code', row.traceability_code)
        
        detectAnomalies(formData).then(res => {
          if (res.code === 200) {
            this.$message.success('异常检测完成')
            // 可以在这里解析返回的检测结果并显示
            if (res.data && res.data !== 'null') {
              const detectionResult = JSON.parse(res.data)
              if (detectionResult.alerts && detectionResult.alerts.length > 0) {
                this.$alert(
                  `检测到 ${detectionResult.alerts.length} 个异常：\n${detectionResult.alerts.map(alert => 
                    `• ${alert.message} (严重等级: ${alert.severity}/5)`
                  ).join('\n')}`,
                  '检测结果',
                  {
                    confirmButtonText: '确定',
                    type: 'warning'
                  }
                )
              } else {
                this.$message.success('未检测到异常')
              }
            }
          } else {
            this.$message.error(res.message || '检测失败')
          }
        }).catch(err => {
          this.$message.error('检测请求失败')
          console.error(err)
        })
      }).catch(() => {
        // 用户取消操作
      })
    },
    
    // 查看报警记录
    handleViewAlerts(row) {
      this.showAlertsDialog = true
      this.alertsLoading = true
      
      const formData = new FormData()
      formData.append('traceability_code', row.traceability_code)
      
      getFruitAlerts(formData).then(res => {
        this.alertsLoading = false
        if (res.code === 200 && res.data) {
          try {
            this.currentAlerts = JSON.parse(res.data)
          } catch (e) {
            this.currentAlerts = []
            console.error('解析报警数据失败:', e)
          }
        } else {
          this.currentAlerts = []
          if (res.message !== '未找到报警记录') {
            this.$message.warning(res.message || '获取报警记录失败')
          }
        }
      }).catch(err => {
        this.alertsLoading = false
        this.$message.error('获取报警记录失败')
        console.error(err)
      })
    },
    
    // 查看召回记录
    handleViewRecalls(row) {
      this.showRecallsDialog = true
      this.recallsLoading = true
      
      const formData = new FormData()
      formData.append('traceability_code', row.traceability_code)
      
      getFruitRecalls(formData).then(res => {
        this.recallsLoading = false
        if (res.code === 200 && res.data) {
          try {
            this.currentRecalls = JSON.parse(res.data)
          } catch (e) {
            this.currentRecalls = []
            console.error('解析召回数据失败:', e)
          }
        } else {
          this.currentRecalls = []
          if (res.message !== '未找到召回记录') {
            this.$message.warning(res.message || '获取召回记录失败')
          }
        }
      }).catch(err => {
        this.recallsLoading = false
        this.$message.error('获取召回记录失败')
        console.error(err)
      })
    },
    
    // 获取严重等级标签类型
    getSeverityType(severity) {
      if (severity >= 4) return 'danger'
      if (severity >= 3) return 'warning'
      if (severity >= 2) return 'info'
      return 'success'
    },
    
    // 获取状态标签类型
    getStatusType(status) {
      switch (status) {
        case 'active': return 'danger'
        case 'resolved': return 'success'
        case 'pending': return 'warning'
        default: return 'info'
      }
    }
  }
}
</script>

<style lang="scss" scoped>
.trace-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  
  .search-section {
    display: flex;
    align-items: center;
  }
}

.demo-table-expand {
  font-size: 0;
}
.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}
.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
}
.trace {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}

.trace-code-cell {
  display: flex;
  align-items: center;
  justify-content: space-between;
  
  .qr-btn {
    margin-left: 10px;
    padding: 2px 8px;
    font-size: 12px;
  }
}

.demo-table-expand {
  font-size: 0;
}

.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}

.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
  display: inline-block;
  vertical-align: top;
}

.demo-table-expand .image-item {
  width: 100%;
  margin-top: 10px;
  margin-bottom: 10px;
}

.demo-table-expand .image-item .el-form-item__content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.alert-detail, .recall-detail {
  line-height: 1.8;
  color: #606266;
}

.alert-detail p, .recall-detail p {
  margin: 8px 0;
}

.alert-detail strong, .recall-detail strong {
  color: #303133;
  margin-right: 8px;
}

</style>
