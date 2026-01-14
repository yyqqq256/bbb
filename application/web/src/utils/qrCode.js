import QRCode from 'qrcode'

/**
 * 生成二维码
 * @param {string} text - 要生成二维码的文本
 * @param {object} options - 二维码配置选项
 * @returns {Promise<string>} - 返回二维码的base64图片数据
 */
export async function generateQRCode(text, options = {}) {
  try {
    const defaultOptions = {
      errorCorrectionLevel: 'H',
      type: 'image/png',
      quality: 0.92,
      margin: 1,
      color: {
        dark: '#000000',
        light: '#FFFFFF'
      },
      width: 256,
      ...options
    }
    
    const qrCodeDataURL = await QRCode.toDataURL(text, defaultOptions)
    return qrCodeDataURL
  } catch (error) {
    console.error('生成二维码失败:', error)
    throw new Error('生成二维码失败')
  }
}

/**
 * 生成产品溯源二维码
 * @param {string} traceabilityCode - 溯源码
 * @param {object} productInfo - 产品信息
 * @returns {Promise<string>} - 返回二维码的base64图片数据
 */
export async function generateProductQRCode(traceabilityCode, productInfo = {}) {
  // 构建二维码内容，包含溯源码和基本信息
  const qrContent = {
    type: 'product_trace',
    code: traceabilityCode,
    name: productInfo.name || '',
    timestamp: new Date().toISOString(),
    version: '1.0'
  }
  
  const qrText = JSON.stringify(qrContent)
  return await generateQRCode(qrText, {
    width: 300,
    margin: 2,
    color: {
      dark: '#096dd9',
      light: '#ffffff'
    }
  })
}

/**
 * 解析二维码内容
 * @param {string} qrContent - 二维码内容
 * @returns {object} - 解析后的数据
 */
export function parseQRCodeContent(qrContent) {
  try {
    const parsed = JSON.parse(qrContent)
    
    // 验证二维码类型
    if (parsed.type === 'product_trace') {
      return {
        isValid: true,
        traceabilityCode: parsed.code,
        productName: parsed.name,
        timestamp: parsed.timestamp,
        version: parsed.version
      }
    }
    
    // 如果不是标准格式，尝试直接作为溯源码处理
    return {
      isValid: true,
      traceabilityCode: qrContent.trim(),
      productName: '',
      timestamp: null,
      version: null
    }
  } catch (error) {
    // JSON解析失败，按纯文本处理
    return {
      isValid: true,
      traceabilityCode: qrContent.trim(),
      productName: '',
      timestamp: null,
      version: null
    }
  }
}

/**
 * 验证二维码内容格式
 * @param {string} qrContent - 二维码内容
 * @returns {boolean} - 是否有效
 */
export function validateQRContent(qrContent) {
  if (!qrContent || typeof qrContent !== 'string') {
    return false
  }
  
  // 检查内容长度
  if (qrContent.length < 3 || qrContent.length > 1000) {
    return false
  }
  
  // 尝试解析
  const parsed = parseQRCodeContent(qrContent)
  return parsed.isValid && parsed.traceabilityCode
}

/**
 * 生成二维码下载链接
 * @param {string} dataURL - base64数据
 * @param {string} filename - 文件名
 * @returns {string} - 下载链接
 */
export function createDownloadLink(dataURL, filename = 'qrcode.png') {
  const link = document.createElement('a')
  link.download = filename
  link.href = dataURL
  return link
}

/**
 * 下载二维码图片
 * @param {string} dataURL - base64数据
 * @param {string} filename - 文件名
 */
export function downloadQRCode(dataURL, filename = 'qrcode.png') {
  const link = createDownloadLink(dataURL, filename)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}