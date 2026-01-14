<template>
  <div class="supply-chain-map">
    <div class="map-header">
      <h3>供应链可视化地图</h3>

      <el-button-group>
        <el-button
          size="small"
          :type="viewMode === 'timeline' ? 'primary' : 'default'"
          @click="viewMode = 'timeline'"
        >
          时间轴视图
        </el-button>

        <el-button
          size="small"
          :type="viewMode === 'network' ? 'primary' : 'default'"
          @click="viewMode = 'network'"
        >
          网络视图
        </el-button>
      </el-button-group>
    </div>

    <div class="map-content">
      <!-- 时间轴视图 -->
      <div v-if="viewMode === 'timeline'" class="timeline-view">
        <el-timeline>
          <el-timeline-item
            v-for="(stage, index) in supplyChainStages"
            :key="index"
            :color="getStageColor(stage.status)"
            :icon="getStageIcon(stage.status)"
            :timestamp="stage.timestamp"
            placement="top"
          >
            <el-card>
              <div slot="header" class="stage-header">
                <span class="stage-title">{{ stage.title }}</span>

                <el-tag
                  :type="getStatusType(stage.status)"
                  size="mini"
                >
                  {{ getStatusText(stage.status) }}
                </el-tag>
              </div>

              <div class="stage-content">
                <div class="stage-info">
                  <p><strong>参与者:</strong> {{ stage.actor }}</p>
                  <p><strong>位置:</strong> {{ stage.location }}</p>
                  <p>
                    <strong>交易ID:</strong>
                    <el-link
                      :href="getBlockchainExplorerUrl(stage.transactionId)"
                      target="_blank"
                      type="primary"
                    >
                      {{ truncateTxId(stage.transactionId) }}
                    </el-link>
                  </p>
                </div>

                <div v-if="stage.details" class="stage-details">
                  <el-collapse>
                    <el-collapse-item name="details" title="详细信息">
                      <el-descriptions :column="2" size="small">
                        <el-descriptions-item
                          v-for="(value, key) in stage.details"
                          :key="key"
                          :label="formatDetailLabel(key)"
                        >
                          {{ value }}
                        </el-descriptions-item>
                      </el-descriptions>
                    </el-collapse-item>
                  </el-collapse>
                </div>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>

      <!-- 网络视图 -->
      <div v-else class="network-view">
        <div ref="networkGraph" class="network-graph">
          <svg class="network-svg" height="400" width="100%">
            <g class="connections">
              <line
                v-for="(connection, index) in networkConnections"
                :key="index"
                :class="['connection-line', connection.status]"
                :x1="connection.x1"
                :x2="connection.x2"
                :y1="connection.y1"
                :y2="connection.y2"
              />
            </g>

            <g class="nodes">
              <g
                v-for="(node, index) in networkNodes"
                :key="index"
                :transform="`translate(${node.x}, ${node.y})`"
                class="network-node"
                @click="showNodeDetails(node)"
              >
                <circle
                  :r="node.radius"
                  :class="['node-circle', node.status]"
                />

                <text
                  class="node-text"
                  text-anchor="middle"
                  :y="5"
                >
                  {{ node.label }}
                </text>

                <text
                  class="node-subtext"
                  text-anchor="middle"
                  :y="25"
                >
                  {{ node.subtext }}
                </text>
              </g>
            </g>
          </svg>
        </div>

        <div class="network-legend">
          <div class="legend-item">
            <span class="legend-color completed" />
            <span>已完成</span>
          </div>
          <div class="legend-item">
            <span class="legend-color in-progress" />
            <span>进行中</span>
          </div>
          <div class="legend-item">
            <span class="legend-color pending" />
            <span>待处理</span>
          </div>
        </div>
      </div>
    </div>

    <el-dialog
      :title="currentNode ? currentNode.label : '节点详情'"
      :visible.sync="nodeDialogVisible"
      width="600px"
    >
      <div v-if="currentNode" class="node-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentNode.status)">
              {{ getStatusText(currentNode.status) }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="位置">
            {{ currentNode.location }}
          </el-descriptions-item>

          <el-descriptions-item label="参与者">
            {{ currentNode.actor }}
          </el-descriptions-item>

          <el-descriptions-item label="时间">
            {{ currentNode.timestamp }}
          </el-descriptions-item>

          <el-descriptions-item :span="2" label="交易ID">
            <el-link
              :href="getBlockchainExplorerUrl(currentNode.transactionId)"
              target="_blank"
              type="primary"
            >
              {{ currentNode.transactionId }}
            </el-link>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentNode.details" class="node-extra-details">
          <h4>详细信息</h4>
          <el-table :data="formatNodeDetails(currentNode.details)" style="width: 100%">
            <el-table-column label="属性" prop="key" width="150" />
            <el-table-column label="值" prop="value" />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'SupplyChainMap',

  props: {
    traceData: {
      type: Object,
      default: () => ({})
    }
  },

  data() {
    return {
      viewMode: 'timeline',
      nodeDialogVisible: false,
      currentNode: null
    }
  },

  computed: {
    supplyChainStages() {
      const stages = []
      const data = this.traceData

      if (data.farmer_input) {
        stages.push({
          title: '原材料/种植',
          actor: data.farmer_input.fa_farmerName,
          location: data.farmer_input.fa_origin,
          timestamp: data.farmer_input.fa_plantTime,
          status: 'completed',
          transactionId: data.farmer_input.fa_txid,
          details: {
            fruitName: data.farmer_input.fa_fruitName,
            pickingTime: data.farmer_input.fa_pickingTime,
            plantTime: data.farmer_input.fa_plantTime
          }
        })
      }

      if (data.factory_input) {
        stages.push({
          title: '生产加工',
          actor: data.factory_input.fac_factoryName,
          location: '工厂地址',
          timestamp: data.factory_input.fac_productionTime,
          status: 'completed',
          transactionId: data.factory_input.fac_txid,
          details: {
            productName: data.factory_input.fac_productName,
            productionBatch: data.factory_input.fac_productionbatch,
            contactNumber: data.factory_input.fac_contactNumber
          }
        })
      }

      if (data.driver_input) {
        stages.push({
          title: '物流运输',
          actor: data.driver_input.dr_name,
          location: '运输途中',
          timestamp: new Date().toISOString(),
          status: 'in-progress',
          transactionId: data.driver_input.dr_txid,
          details: {
            driverName: data.driver_input.dr_name,
            carNumber: data.driver_input.dr_carNumber,
            phone: data.driver_input.dr_phone,
            transport: data.driver_input.dr_transport
          }
        })
      }

      if (data.shop_input) {
        stages.push({
          title: '商店销售',
          actor: data.shop_input.sh_shopName,
          location: data.shop_input.sh_shopAddress,
          timestamp: data.shop_input.sh_storeTime,
          status: 'pending',
          transactionId: data.shop_input.sh_txid,
          details: {
            storeTime: data.shop_input.sh_storeTime,
            sellTime: data.shop_input.sh_sellTime,
            shopPhone: data.shop_input.sh_shopPhone
          }
        })
      }

      return stages
    },

    networkNodes() {
      const nodes = []
      const centerX = 200
      const centerY = 200
      const radius = 150

      this.supplyChainStages.forEach((stage, index) => {
        const angle = (index * 2 * Math.PI) / this.supplyChainStages.length - Math.PI / 2
        const x = centerX + radius * Math.cos(angle)
        const y = centerY + radius * Math.sin(angle)

        nodes.push({
          ...stage,
          x,
          y,
          radius: 40,
          label: this.getShortTitle(stage.title),
          subtext: this.truncateText(stage.actor, 8)
        })
      })

      return nodes
    },

    networkConnections() {
      const connections = []

      for (let i = 0; i < this.networkNodes.length - 1; i++) {
        const current = this.networkNodes[i]
        const next = this.networkNodes[i + 1]

        connections.push({
          x1: current.x,
          y1: current.y,
          x2: next.x,
          y2: next.y,
          status: current.status
        })
      }

      return connections
    }
  },

  methods: {
    getStageColor(status) {
      const colors = {
        completed: '#67c23a',
        'in-progress': '#409eff',
        pending: '#e6a23c',
        failed: '#f56c6c'
      }

      return colors[status] || '#909399'
    },

    getStageIcon(status) {
      const icons = {
        completed: 'el-icon-check',
        'in-progress': 'el-icon-loading',
        pending: 'el-icon-time',
        failed: 'el-icon-close'
      }

      return icons[status] || 'el-icon-info'
    },

    getStatusType(status) {
      const types = {
        completed: 'success',
        'in-progress': 'primary',
        pending: 'warning',
        failed: 'danger'
      }

      return types[status] || 'info'
    },

    getStatusText(status) {
      const texts = {
        completed: '已完成',
        'in-progress': '进行中',
        pending: '待处理',
        failed: '失败'
      }

      return texts[status] || '未知'
    },

    getBlockchainExplorerUrl(txId) {
      return `https://fabric-explorer.example.com/tx/${txId}`
    },

    truncateTxId(txId) {
      if (!txId) return 'N/A'

      return txId.length > 16
        ? `${txId.substring(0, 8)}...${txId.substring(txId.length - 8)}`
        : txId
    },

    formatDetailLabel(key) {
      const labels = {
        fruitName: '产品名称',
        pickingTime: '采摘时间',
        plantTime: '种植时间',
        productName: '商品名称',
        productionBatch: '生产批次',
        contactNumber: '联系电话',
        driverName: '司机姓名',
        carNumber: '车牌号',
        phone: '电话',
        transport: '运输记录',
        storeTime: '入库时间',
        sellTime: '销售时间',
        shopPhone: '商店电话'
      }

      return labels[key] || key
    },

    getShortTitle(title) {
      const shortTitles = {
        '原材料/种植': '种植',
        '生产加工': '生产',
        '物流运输': '物流',
        '商店销售': '销售'
      }

      return shortTitles[title] || title.substring(0, 2)
    },

    truncateText(text, maxLength) {
      if (!text) return ''

      return text.length > maxLength
        ? `${text.substring(0, maxLength)}...`
        : text
    },

    showNodeDetails(node) {
      this.currentNode = node
      this.nodeDialogVisible = true
    },

    formatNodeDetails(details) {
      return Object.entries(details).map(([key, value]) => ({
        key: this.formatDetailLabel(key),
        value
      }))
    }
  }
}
</script>

<style lang="scss" scoped>
/* 样式结构保持原意，仅修复空格/大小写/对齐规范，功能无变动 */
