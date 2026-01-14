<template>
  <div class="product-qr-display">
    <el-dialog
      :visible.sync="visible"
      :title="`产品二维码 - ${productName}`"
      width="400px"
      @close="handleClose"
    >
      <div class="qr-content" v-loading="generating">
        <div v-if="qrCodeData" class="qr-image-section">
          <img :src="qrCodeData" alt="产品二维码" class="qr-image" />
          <div class="qr-info">
            <p class="trace-code">溯源码: {{ traceabilityCode }}</p>
            <p class="generate-time">生成时间: {{ generateTime }}</p>
          </div>
        </div>
        
        <div v-if="error" class="error-message">
          <el-alert
            :title="error"
            type="error"
            show-icon
            :closable="false"
          />
        </div>
      </div>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button 
          type="primary" 
          @click="downloadQRCode"
          :disabled="!qrCodeData"
          icon="el-icon-download"
        >
          下载二维码
        </el-button>
        <el-button 
          type="success" 
          @click="printQRCode"
          :disabled="!qrCodeData"
          icon="el-icon-printer"
        >
          打印
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { generateProductQRCode, downloadQRCode as downloadQR } from '@/utils/qrCode'

export default {
  name: 'ProductQRDisplay',
  props: {
    visible: {
      type: Boolean,
      default: false
    },
    traceabilityCode: {
      type: String,
      required: true
    },
    productName: {
      type: String,
      default: ''
    },
    productInfo: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      qrCodeData: null,
      generating: false,
      error: null,
      generateTime: ''
    }
  },
  watch: {
    visible: {
      immediate: true,
      handler(val) {
        if (val) {
          this.generateQRCode()
        }
      }
    }
  },
  methods: {
    async generateQRCode() {
      if (!this.traceabilityCode) {
        this.error = '溯源码不能为空'
        return
      }
      
      this.generating = true
      this.error = null
      
      try {
        const qrData = await generateProductQRCode(
          this.traceabilityCode,
          {
            name: this.productName,
            ...this.productInfo
          }
        )
        
        this.qrCodeData = qrData
        this.generateTime = new Date().toLocaleString()
      } catch (error) {
        console.error('生成二维码失败:', error)
        this.error = '生成二维码失败，请稍后重试'
      } finally {
        this.generating = false
      }
    },
    
    downloadQRCode() {
      if (this.qrCodeData) {
        const filename = `产品二维码_${this.traceabilityCode}_${Date.now()}.png`
        downloadQR(this.qrCodeData, filename)
        this.$message.success('二维码下载成功')
      }
    },
    
    printQRCode() {
      if (this.qrCodeData) {
        const printWindow = window.open('', '_blank', 'width=400,height=500')
        printWindow.document.write(`
          <html>
            <head>
              <title>产品二维码 - ${this.productName}</title>
              <style>
                body { 
                  font-family: Arial, sans-serif; 
                  text-align: center; 
                  padding: 20px; 
                }
                .qr-container { 
                  margin: 20px auto; 
                }
                .qr-image { 
                  max-width: 300px; 
                  border: 1px solid #ccc; 
                  padding: 10px; 
                }
                .product-info { 
                  margin-top: 15px; 
                  font-size: 14px; 
                }
                .trace-code { 
                  font-weight: bold; 
                  color: #333; 
                  margin: 10px 0; 
                }
              </style>
            </head>
            <body>
              <div class="qr-container">
                <h2>${this.productName || '产品二维码'}</h2>
                <img src="${this.qrCodeData}" alt="产品二维码" class="qr-image" />
                <div class="product-info">
                  <p class="trace-code">溯源码: ${this.traceabilityCode}</p>
                  <p>生成时间: ${this.generateTime}</p>
                </div>
              </div>
            </body>
          </html>
        `)
        printWindow.document.close()
        printWindow.print()
      }
    },
    
    handleClose() {
      this.qrCodeData = null
      this.error = null
      this.generateTime = ''
      this.$emit('update:visible', false)
    }
  }
}
</script>

<style lang="scss" scoped>
.product-qr-display {
  .qr-content {
    text-align: center;
    min-height: 200px;
    
    .qr-image-section {
      .qr-image {
        max-width: 300px;
        border: 1px solid #dcdfe6;
        border-radius: 8px;
        padding: 10px;
        background: white;
      }
      
      .qr-info {
        margin-top: 15px;
        padding: 10px;
        background: #f5f7fa;
        border-radius: 4px;
        
        p {
          margin: 5px 0;
          font-size: 14px;
          color: #606266;
        }
        
        .trace-code {
          font-weight: bold;
          color: #303133;
        }
      }
    }
    
    .error-message {
      padding: 20px;
    }
  }
  
  .dialog-footer {
    text-align: right;
  }
}
</style>