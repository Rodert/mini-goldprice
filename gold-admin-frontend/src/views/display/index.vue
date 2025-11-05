<template>
  <div class="gold-display">
    <!-- 顶部横幅 -->
    <div class="header">
      <div class="logo">
        <span class="logo-icon">💎</span>
        <span class="brand-name">SineGem 中国珠宝</span>
      </div>
      <div class="title">今日金价</div>
    </div>

    <!-- 主体内容 -->
    <div class="main-content">
      <!-- 左侧：视频 -->
      <div class="video-section">
        <video
          ref="videoPlayer"
          :src="videoUrl"
          autoplay
          loop
          muted
          playsinline
          @timeupdate="updateProgress"
          @loadedmetadata="onVideoLoaded"
        ></video>
        <div v-if="!videoUrl" class="video-placeholder">
          <i class="el-icon-video-camera"></i>
          <p>视频展示区域</p>
          <p class="hint">请在 public/assets 目录下添加 jewelry-video.mp4</p>
        </div>
      </div>

      <!-- 中间：产品图片滚动 -->
      <div class="scroll-images">
        <div class="scroll-container" :style="{ transform: `translateY(${scrollY}px)` }">
          <div v-for="(img, index) in doubleImages" :key="index" class="product-image">
            <img :src="img" alt="产品图片" />
          </div>
        </div>
      </div>

      <!-- 右侧：金价表格 -->
      <div class="price-table">
        <!-- 调试信息 -->
        <div v-if="priceList.length === 0" style="color: red; padding: 20px; text-align: center;">
          数据加载中... 当前数据条数: {{ priceList.length }}
        </div>
        
        <table v-if="priceList.length > 0">
          <thead>
            <tr>
              <th>品名</th>
              <th>销售价</th>
              <th>工费</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in priceList" :key="item.id || item.name">
              <td class="product-name">{{ item.name }}</td>
              <td class="price">{{ item.price }}</td>
              <td class="fee">{{ item.fee }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 底部信息栏 -->
    <div class="footer">
      <div class="datetime">{{ currentTime }}</div>
      <div class="marquee">
        <div class="marquee-content">
          <div class="marquee-text" :style="{ transform: `translateX(${marqueeX}px)` }">
            中国珠宝欢迎您！！！
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getPriceList } from '@/api/price'

export default {
  name: 'GoldDisplay',
  data() {
    return {
      // 视频
      videoUrl: 'https://media.w3.org/2010/05/sintel/trailer.mp4',
      currentVideoTime: 0,
      videoDuration: 0,
      videoProgress: 0,
      
      // 产品图片（垂直滚动）
      productImages: [
        'https://qcloud.dpfile.com/pc/YCTUsa0Z4mrzK7qBdGKcwTKtKkDKaHEnpRurI7Y593BGwRK899Q_dG3QMjFCppQBY0q73sB2DyQcgmKUxZFQtw.jpg',
        'https://qcloud.dpfile.com/pc/YCTUsa0Z4mrzK7qBdGKcwTKtKkDKaHEnpRurI7Y593BGwRK899Q_dG3QMjFCppQBY0q73sB2DyQcgmKUxZFQtw.jpg',
        'https://qcloud.dpfile.com/pc/YCTUsa0Z4mrzK7qBdGKcwTKtKkDKaHEnpRurI7Y593BGwRK899Q_dG3QMjFCppQBY0q73sB2DyQcgmKUxZFQtw.jpg',
        'https://qcloud.dpfile.com/pc/YCTUsa0Z4mrzK7qBdGKcwTKtKkDKaHEnpRurI7Y593BGwRK899Q_dG3QMjFCppQBY0q73sB2DyQcgmKUxZFQtw.jpg',
        'https://qcloud.dpfile.com/pc/YCTUsa0Z4mrzK7qBdGKcwTKtKkDKaHEnpRurI7Y593BGwRK899Q_dG3QMjFCppQBY0q73sB2DyQcgmKUxZFQtw.jpg'
      ],
      scrollY: 0,
      scrollSpeed: 30, // 滚动速度（像素/秒）
      
      // 价格表
      priceList: [],
      
      // 时间
      currentTime: '',
      
      // 跑马灯
      marqueeX: 0,
      marqueeWidth: 0,
      marqueeSpeed: 100, // 滚动速度（像素/秒）
      
      // 动画定时器
      scrollTimer: null,
      marqueeTimer: null,
      timeTimer: null,
      refreshTimer: null
    }
  },
  computed: {
    // 双倍图片数组用于无缝滚动
    doubleImages() {
      return [...this.productImages, ...this.productImages]
    }
  },
  mounted() {
    console.log('🚀 页面加载完成，开始初始化...')
    console.log('📊 初始 priceList:', this.priceList)
    
    this.fetchPriceData()
    this.startAnimations()
    this.updateTime()
    this.initMarquee()
    
    // 定期刷新价格数据（每5分钟）
    this.refreshTimer = setInterval(() => {
      this.fetchPriceData()
    }, 5 * 60 * 1000)
    
    // 5秒后再次检查数据
    setTimeout(() => {
      console.log('⏰ 5秒后检查 priceList:', this.priceList)
    }, 5000)
  },
  beforeDestroy() {
    this.stopAnimations()
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer)
    }
  },
  methods: {
    // 获取金价数据
    async fetchPriceData() {
      try {
        console.log('🔍 开始获取金价数据...')
        const response = await getPriceList({
          page: 1,
          page_size: 100
        })
        
        console.log('📊 完整API响应:', response)
        console.log('📊 response.code:', response.code)
        console.log('📊 response.data:', response.data)
        
        // 检查多种可能的数据结构
        if (response && response.code === 200) {
          let dataList = null
          
          // 尝试不同的数据结构
          if (response.data && response.data.list) {
            dataList = response.data.list
            console.log('✅ 使用 response.data.list')
          } else if (response.data && Array.isArray(response.data)) {
            dataList = response.data
            console.log('✅ 使用 response.data (数组)')
          } else if (Array.isArray(response)) {
            dataList = response
            console.log('✅ 使用 response (数组)')
          }
          
          if (dataList && dataList.length > 0) {
            console.log('📊 原始数据列表:', dataList)
            this.priceList = dataList.slice(0, 5).map(item => ({
              id: item.id,
              name: item.name,
              price: item.sell_price || item.price || 0,  // 使用销售价
              fee: item.fee || 10
            }))
            console.log('✅ 金价数据加载成功，条数:', this.priceList.length)
            console.log('✅ 处理后的数据:', JSON.stringify(this.priceList, null, 2))
            
            // 强制更新视图
            this.$forceUpdate()
          } else {
            console.warn('⚠️ 数据列表为空，使用模拟数据')
            this.useMockData()
          }
        } else {
          console.warn('⚠️ API返回格式不正确，使用模拟数据')
          console.log('⚠️ response.code =', response?.code)
          this.useMockData()
        }
      } catch (error) {
        console.error('❌ 获取价格数据失败:', error)
        this.useMockData()
      }
    },
    
    // 使用模拟数据
    useMockData() {
      console.log('📝 使用模拟数据')
      this.priceList = [
        { id: 1, name: '足金(9999)', price: 488, fee: 10, image: '' },
        { id: 2, name: '足金(999)', price: 428, fee: 10, image: '' },
        { id: 3, name: 'Pt950', price: 388, fee: 10, image: '' },
        { id: 4, name: 'Pt990', price: 408, fee: 10, image: '' },
        { id: 5, name: 'PD950', price: 218, fee: 10, image: '' }
      ]
      console.log('✅ 模拟数据已设置:', this.priceList)
    },
    
    // 初始化跑马灯
    initMarquee() {
      this.$nextTick(() => {
        const marqueeElement = this.$el.querySelector('.marquee-text')
        if (marqueeElement) {
          this.marqueeWidth = marqueeElement.offsetWidth
          this.marqueeX = this.$el.querySelector('.marquee-content').offsetWidth
        }
      })
    },
    
    // 开始所有动画
    startAnimations() {
      // 产品图片垂直滚动
      this.scrollTimer = setInterval(() => {
        this.scrollY -= this.scrollSpeed / 60
        const singleHeight = this.productImages.length * 220 // 每张图200px + 20px padding
        if (Math.abs(this.scrollY) >= singleHeight) {
          this.scrollY = 0
        }
      }, 1000 / 60)
      
      // 跑马灯横向滚动
      this.marqueeTimer = setInterval(() => {
        this.marqueeX -= this.marqueeSpeed / 60
        const containerWidth = this.$el.querySelector('.marquee-content')?.offsetWidth || 500
        if (this.marqueeX < -this.marqueeWidth) {
          this.marqueeX = containerWidth
        }
      }, 1000 / 60)
      
      // 更新时间
      this.timeTimer = setInterval(() => {
        this.updateTime()
      }, 1000)
    },
    
    // 停止所有动画
    stopAnimations() {
      if (this.scrollTimer) clearInterval(this.scrollTimer)
      if (this.marqueeTimer) clearInterval(this.marqueeTimer)
      if (this.timeTimer) clearInterval(this.timeTimer)
    },
    
    // 更新时间显示
    updateTime() {
      const now = new Date()
      const days = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六']
      const year = now.getFullYear()
      const month = String(now.getMonth() + 1).padStart(2, '0')
      const date = String(now.getDate()).padStart(2, '0')
      const day = days[now.getDay()]
      const hours = String(now.getHours()).padStart(2, '0')
      const minutes = String(now.getMinutes()).padStart(2, '0')
      const seconds = String(now.getSeconds()).padStart(2, '0')
      
      this.currentTime = `${year}-${month}-${date} ${day} ${hours}:${minutes}:${seconds}`
    },
    
    // 视频加载完成
    onVideoLoaded() {
      const video = this.$refs.videoPlayer
      if (video) {
        this.videoDuration = video.duration || 0
      }
    },
    
    // 更新视频进度
    updateProgress() {
      const video = this.$refs.videoPlayer
      if (video) {
        this.currentVideoTime = video.currentTime
        this.videoDuration = video.duration || 0
        this.videoProgress = this.videoDuration > 0 
          ? (this.currentVideoTime / this.videoDuration) * 100 
          : 0
      }
    },
    
    // 格式化时间（秒转为 mm:ss）
    formatTime(seconds) {
      if (!seconds || isNaN(seconds)) return '00:00'
      const mins = Math.floor(seconds / 60)
      const secs = Math.floor(seconds % 60)
      return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`
    }
  }
}
</script>

<style lang="scss" scoped>
.gold-display {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  color: #f0c674;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  font-family: 'Microsoft YaHei', 'PingFang SC', Arial, sans-serif;
}

// 顶部横幅
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 40px;
  background: rgba(0, 0, 0, 0.5);
  border-bottom: 2px solid #f0c674;
  box-shadow: 0 2px 10px rgba(240, 198, 116, 0.3);
  
  .logo {
    display: flex;
    align-items: center;
    gap: 15px;
    
    .logo-icon {
      font-size: 48px;
      filter: drop-shadow(0 0 10px rgba(240, 198, 116, 0.6));
    }
    
    .brand-name {
      font-size: 36px;
      font-weight: bold;
      letter-spacing: 2px;
      text-shadow: 0 0 20px rgba(240, 198, 116, 0.5);
    }
  }
  
  .title {
    font-size: 42px;
    font-weight: bold;
    letter-spacing: 4px;
    text-shadow: 0 0 20px rgba(240, 198, 116, 0.5);
  }
}

// 主体内容
.main-content {
  flex: 1;
  display: flex;
  gap: 20px;
  padding: 20px;
  overflow: hidden;
}

// 视频区域（60%宽度，保持16:9比例）
.video-section {
  flex: 6;
  background: #000;
  border-radius: 10px;
  overflow: hidden;
  position: relative;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  
  video {
    width: 100%;
    height: 100%;
    object-fit: contain;  // 保持视频16:9比例
  }
  
  .video-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    color: #666;
    
    i {
      font-size: 80px;
      margin-bottom: 20px;
    }
    
    p {
      font-size: 20px;
      margin: 5px 0;
    }
    
    .hint {
      font-size: 14px;
      color: #888;
    }
  }
}

// 产品图片滚动区域（10%宽度，窄条）
.scroll-images {
  flex: 1;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 10px;
  overflow: hidden;
  position: relative;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  
  .scroll-container {
    .product-image {
      height: 200px;
      padding: 10px;
      
      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 8px;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
      }
    }
  }
}

// 价格表格区域（30%宽度）
.price-table {
  flex: 3;
  background: rgba(0, 0, 0, 0.5);
  border-radius: 10px;
  padding: 20px;
  overflow: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  
  table {
    width: 100%;
    border-collapse: collapse;
    
    thead tr {
      border-bottom: 2px solid rgba(240, 198, 116, 0.5);
    }
    
    tbody tr {
      min-height: 80px;
      max-height: 100px;
      transition: background-color 0.3s;
      border-bottom: 1px solid rgba(240, 198, 116, 0.2);
      
      &:hover {
        background: rgba(240, 198, 116, 0.05);
      }
      
      &:last-child {
        border-bottom: none;
      }
    }
    
    th, td {
      text-align: center;
      padding: 15px 10px;
      vertical-align: middle;
    }
    
    th {
      font-size: 28px;
      font-weight: bold;
      background: rgba(240, 198, 116, 0.1);
      color: #f0c674;
      text-shadow: 0 0 10px rgba(240, 198, 116, 0.3);
      padding: 20px 10px;
    }
    
    td {
      font-size: 32px;
      
      &.product-name {
        font-weight: bold;
        font-size: 30px;
      }
      
      &.price {
        font-size: 36px;
        color: #ff6b6b;
        font-weight: bold;
        text-shadow: 0 0 10px rgba(255, 107, 107, 0.5);
      }
      
      &.fee {
        font-size: 28px;
        color: #87ceeb;
      }
    }
  }
}

// 底部信息栏
.footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 40px;
  background: rgba(0, 0, 0, 0.5);
  border-top: 2px solid #f0c674;
  box-shadow: 0 -2px 10px rgba(240, 198, 116, 0.3);
  
  .datetime {
    flex: 1;
    font-size: 20px;
    font-weight: 500;
  }
  
  .marquee {
    flex: 1;
    overflow: hidden;
    position: relative;
    
    .marquee-content {
      width: 100%;
      height: 30px;
      position: relative;
      overflow: hidden;
    }
    
    .marquee-text {
      font-size: 24px;
      font-weight: bold;
      white-space: nowrap;
      position: absolute;
      left: 0;
      top: 0;
      text-shadow: 0 0 10px rgba(240, 198, 116, 0.5);
    }
  }
}

// 滚动条样式
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: rgba(240, 198, 116, 0.5);
  border-radius: 4px;
  
  &:hover {
    background: rgba(240, 198, 116, 0.7);
  }
}
</style>

