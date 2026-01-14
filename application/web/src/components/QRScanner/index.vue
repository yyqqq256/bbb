<template>
  <div class="qr-scanner-container">
    <div class="scanner-header">
      <h3>二维码扫描</h3>

      <el-button
        class="close-btn"
        icon="el-icon-close"
        type="text"
        @click="$emit('close')"
      />
    </div>

    <div class="scanner-content">
      <div v-if="!scanning" class="scanner-setup">
        <el-button
          icon="el-icon-camera"
          type="primary"
          @click="startScanning"
        >
          开始扫描
        </el-button>
      </div>

      <div v-else class="scanner-active">
        <div id="qr-reader" class="qr-reader" />
        <el-button
          icon="el-icon-video-pause"
          type="warning"
          @click="stopScanning"
        >
          停止扫描
        </el-button>
      </div>

      <div v-if="scanResult" class="scan-result">
        <el-alert
          :description="scanResult.message"
          :title="scanResult.success ? '扫描成功' : '扫描失败'"
          :type="scanResult.success ? 'success' : 'error'"
          show-icon
          @close="clearResult"
        />
      </div>
    </div>

    <div class="manual-input-section">
      <el-divider>或手动输入</el-divider>

      <el-input
        v-model="manualCode"
        clearable
        placeholder="请输入溯源码"
        @keyup.enter="handleManualInput"
      >
        <el-button
          slot="append"
          type="primary"
          @click="handleManualInput"
        >
          查询
        </el-button>
      </el-input>
    </div>
  </div>
</template>

<script>
import { Html5Qrcode } from 'html5-qrcode'

export default {
  name: 'QRScanner',

  data() {
    return {
      scanning: false,
      html5QrCode: null,
      scanResult: null,
      manualCode: ''
    }
  },

  beforeDestroy() {
    if (this.html5QrCode) {
      this.stopScanning()
    }
  },

  methods: {
    async startScanning() {
      try {
        this.html5QrCode = new Html5Qrcode('qr-reader')

        const config = {
          fps: 10,
          qrbox: { width: 250, height: 250 },
          aspectRatio: 1.0
        }

        await this.html5QrCode.start(
          { facingMode: 'environment' },
          config,
          this.onScanSuccess,
          this.onScanFailure
        )

        this.scanning = true
        this.$message.success('扫描已启动')
      } catch (err) {
        console.error('启动扫描失败:', err)
        this.$message.error('无法启动摄像头，请检查权限设置')
      }
    },

    async stopScanning() {
      if (!this.html5QrCode) return

      try {
        await this.html5QrCode.stop()
        this.html5QrCode.clear()
        this.scanning = false
        this.$message.info('扫描已停止')
      } catch (err) {
        console.error('停止扫描失败:', err)
      }
    },

    onScanSuccess(decodedText) {
      console.log('二维码扫描成功:', decodedText)

      this.scanResult = {
        success: true,
        message: `扫描结果: ${decodedText}`
      }

      this.stopScanning()
      this.$emit('scan-success', decodedText)
    },

    onScanFailure(err) {
      console.log('扫描中...', err)
    },

    handleManualInput() {
      if (!this.manualCode.trim()) {
        this.$message.warning('请输入溯源码')
        return
      }

      this.$emit('scan-success', this.manualCode.trim())
    },

    clearResult() {
      this.scanResult = null
    }
  }
}
</script>

<style lang="scss" scoped>
.qr-scanner-container {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  max-width: 500px;
  margin: 0 auto;
}

.scanner-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h3 {
    margin: 0;
    color: #303133;
  }

  .close-btn {
    padding: 0;
    font-size: 18px;
  }
}

.scanner-content {
  text-align: center;
  margin-bottom: 20px;
}

.scanner-setup {
  padding: 40px 0;
}

.scanner-active {
  .qr-reader {
    width: 100%;
    max-width: 400px;
    height: 400px;
    margin: 0 auto 20px;
    border: 2px solid #409eff;
    border-radius: 8px;
    overflow: hidden;
  }
}

.scan-result {
  margin-top: 20px;
}

.manual-input-section {
  margin-top: 30px;
}
</style>
