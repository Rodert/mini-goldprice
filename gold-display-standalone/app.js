/**
 * 金价展示页面 - 原生JavaScript实现
 * 兼容旧版浏览器（Android 4.4+, Chrome 30+）
 */

(function() {
    'use strict';

    // 配置项
    var CONFIG = {
        // API配置
        apiBaseUrl: 'http://gold.javapub.net.cn/api/openapi',
        apiTimeout: 10000, // 10秒超时
        
        // 动画配置
        scrollSpeed: 30,      // 产品图片滚动速度（像素/秒）
        marqueeSpeed: 100,    // 跑马灯滚动速度（像素/秒）
        animationFPS: 60,    // 动画帧率
        
        // 刷新配置
        refreshInterval: 5 * 60 * 1000, // 5分钟刷新一次
        
        // 资源路径
        videoPath: 'assets/jewelry-video.mp4',
        productImages: [
            'assets/jiezhi.png',
            'assets/jiezhi.png',
            'assets/shouzhuo.png',
            'assets/jiezhi.png',
            'assets/shouzhuo.png',
        ],
        
        // 模拟数据（API失败时使用）
        mockData: [
            { id: 1, name: '足金(9999)', price: 488, fee: 10 },
            { id: 2, name: '足金(999)', price: 428, fee: 10 },
            { id: 3, name: 'Pt950', price: 388, fee: 10 },
            { id: 4, name: 'Pt990', price: 408, fee: 10 },
            { id: 5, name: 'PD950', price: 218, fee: 10 }
        ]
    };

    // 全局变量
    var state = {
        priceList: [],
        scrollY: 0,
        marqueeX: 0,
        marqueeWidth: 0,
        scrollTimer: null,
        marqueeTimer: null,
        timeTimer: null,
        refreshTimer: null
    };

    // DOM元素引用
    var elements = {
        videoPlayer: null,
        videoPlaceholder: null,
        scrollContainer: null,
        priceTable: null,
        priceTableBody: null,
        loadingMessage: null,
        datetime: null,
        marqueeText: null
    };

    /**
     * 初始化 - 页面加载完成后执行
     */
    function init() {
        console.log('🚀 页面初始化开始...');
        
        // 获取DOM元素
        getElements();
        
        // 初始化视频
        initVideo();
        
        // 初始化产品图片
        initProductImages();
        
        // 初始化跑马灯
        initMarquee();
        
        // 获取金价数据
        fetchPriceData();
        
        // 启动动画
        startAnimations();
        
        // 更新时间
        updateTime();
        
        // 定期刷新
        startRefreshTimer();
        
        console.log('✅ 页面初始化完成');
    }

    /**
     * 获取DOM元素
     */
    function getElements() {
        elements.videoPlayer = document.getElementById('videoPlayer');
        elements.videoPlaceholder = document.getElementById('videoPlaceholder');
        elements.scrollContainer = document.getElementById('scrollContainer');
        elements.priceTable = document.getElementById('priceTable');
        elements.priceTableBody = document.getElementById('priceTableBody');
        elements.loadingMessage = document.getElementById('loadingMessage');
        elements.datetime = document.getElementById('datetime');
        elements.marqueeText = document.getElementById('marqueeText');
    }

    /**
     * 初始化视频
     */
    function initVideo() {
        if (!elements.videoPlayer) return;
        
        var video = elements.videoPlayer;
        
        // 检查视频是否加载成功
        video.addEventListener('error', function() {
            console.warn('⚠️ 视频加载失败，显示占位符');
            if (elements.videoPlaceholder) {
                elements.videoPlaceholder.style.display = 'flex';
            }
            if (video) {
                video.style.display = 'none';
            }
        });
        
        video.addEventListener('loadeddata', function() {
            console.log('✅ 视频加载成功');
            if (elements.videoPlaceholder) {
                elements.videoPlaceholder.style.display = 'none';
            }
        });
        
        // 设置视频源
        video.src = CONFIG.videoPath;
    }

    /**
     * 初始化产品图片
     */
    function initProductImages() {
        if (!elements.scrollContainer) return;
        
        var images = CONFIG.productImages;
        var doubleImages = images.concat(images); // 复制数组实现无缝滚动
        
        // 清空容器
        elements.scrollContainer.innerHTML = '';
        
        // 创建图片元素
        for (var i = 0; i < doubleImages.length; i++) {
            var imgDiv = document.createElement('div');
            imgDiv.className = 'product-image';
            
            var img = document.createElement('img');
            img.src = doubleImages[i];
            img.alt = '产品图片 ' + (i + 1);
            
            // 图片加载错误处理
            img.onerror = function() {
                this.style.display = 'none';
            };
            
            imgDiv.appendChild(img);
            elements.scrollContainer.appendChild(imgDiv);
        }
    }

    /**
     * 初始化跑马灯
     */
    function initMarquee() {
        if (!elements.marqueeText) return;
        
        // 等待DOM渲染完成后获取宽度
        setTimeout(function() {
            if (elements.marqueeText) {
                state.marqueeWidth = elements.marqueeText.offsetWidth || 500;
                
                // 获取容器宽度
                var marqueeContent = elements.marqueeText.parentElement;
                if (marqueeContent) {
                    state.marqueeX = marqueeContent.offsetWidth || 500;
                }
            }
        }, 100);
    }

    /**
     * 获取金价数据
     */
    function fetchPriceData() {
        console.log('🔍 开始获取金价数据...');
        
        var url = CONFIG.apiBaseUrl + '/prices?page=1&page_size=100';
        
        // 使用XMLHttpRequest（兼容旧版浏览器）
        var xhr = new XMLHttpRequest();
        var timeout = setTimeout(function() {
            xhr.abort();
            console.warn('⚠️ API请求超时，使用模拟数据');
            useMockData();
        }, CONFIG.apiTimeout);
        
        xhr.open('GET', url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        
        xhr.onreadystatechange = function() {
            if (xhr.readyState === 4) {
                clearTimeout(timeout);
                
                if (xhr.status === 200) {
                    try {
                        var response = JSON.parse(xhr.responseText);
                        handlePriceResponse(response);
                    } catch (e) {
                        console.error('❌ 解析响应数据失败:', e);
                        useMockData();
                    }
                } else {
                    console.warn('⚠️ API请求失败，状态码:', xhr.status);
                    useMockData();
                }
            }
        };
        
        xhr.onerror = function() {
            clearTimeout(timeout);
            console.error('❌ 网络错误，使用模拟数据');
            useMockData();
        };
        
        xhr.send();
    }

    /**
     * 处理价格响应数据
     */
    function handlePriceResponse(response) {
        console.log('📊 API响应:', response);
        
        var dataList = null;
        
        // 尝试不同的数据结构
        if (response && response.code === 200) {
            if (response.data && response.data.list) {
                dataList = response.data.list;
                console.log('✅ 使用 response.data.list');
            } else if (response.data && Array.isArray(response.data)) {
                dataList = response.data;
                console.log('✅ 使用 response.data (数组)');
            } else if (Array.isArray(response)) {
                dataList = response;
                console.log('✅ 使用 response (数组)');
            }
        }
        
        if (dataList && dataList.length > 0) {
            // 取前5条数据
            var displayData = dataList.slice(0, 5);
            
            state.priceList = displayData.map(function(item) {
                return {
                    id: item.id || 0,
                    name: item.name || '',
                    price: item.sell_price || item.price || 0,
                    fee: item.fee || 10
                };
            });
            
            console.log('✅ 金价数据加载成功，条数:', state.priceList.length);
            renderPriceTable();
        } else {
            console.warn('⚠️ 数据列表为空，使用模拟数据');
            useMockData();
        }
    }

    /**
     * 使用模拟数据
     */
    function useMockData() {
        console.log('📝 使用模拟数据');
        state.priceList = CONFIG.mockData.slice();
        renderPriceTable();
    }

    /**
     * 渲染价格表格
     */
    function renderPriceTable() {
        if (!elements.priceTable || !elements.priceTableBody) return;
        
        // 隐藏加载消息
        if (elements.loadingMessage) {
            elements.loadingMessage.style.display = 'none';
        }
        
        // 清空表格
        elements.priceTableBody.innerHTML = '';
        
        // 渲染数据
        for (var i = 0; i < state.priceList.length; i++) {
            var item = state.priceList[i];
            var row = document.createElement('tr');
            
            var nameCell = document.createElement('td');
            nameCell.className = 'product-name';
            nameCell.textContent = item.name;
            
            var priceCell = document.createElement('td');
            priceCell.className = 'price';
            priceCell.textContent = item.price;
            
            var feeCell = document.createElement('td');
            feeCell.className = 'fee';
            feeCell.textContent = item.fee;
            
            row.appendChild(nameCell);
            row.appendChild(priceCell);
            row.appendChild(feeCell);
            
            elements.priceTableBody.appendChild(row);
        }
        
        // 显示表格
        elements.priceTable.style.display = 'table';
        
        console.log('✅ 价格表格渲染完成');
    }

    /**
     * 启动所有动画
     */
    function startAnimations() {
        var frameTime = 1000 / CONFIG.animationFPS;
        
        // 产品图片垂直滚动
        state.scrollTimer = setInterval(function() {
            if (!elements.scrollContainer) return;
            
            state.scrollY -= CONFIG.scrollSpeed / CONFIG.animationFPS;
            
            // 计算单组图片高度
            var singleHeight = CONFIG.productImages.length * 220; // 200px + 20px padding
            
            // 重置位置实现无缝循环
            if (Math.abs(state.scrollY) >= singleHeight) {
                state.scrollY = 0;
            }
            
            // 应用transform（兼容旧版浏览器）
            var transform = 'translateY(' + state.scrollY + 'px)';
            elements.scrollContainer.style.webkitTransform = transform;
            elements.scrollContainer.style.mozTransform = transform;
            elements.scrollContainer.style.msTransform = transform;
            elements.scrollContainer.style.oTransform = transform;
            elements.scrollContainer.style.transform = transform;
        }, frameTime);
        
        // 跑马灯横向滚动
        state.marqueeTimer = setInterval(function() {
            if (!elements.marqueeText) return;
            
            state.marqueeX -= CONFIG.marqueeSpeed / CONFIG.animationFPS;
            
            // 获取容器宽度
            var marqueeContent = elements.marqueeText.parentElement;
            var containerWidth = marqueeContent ? marqueeContent.offsetWidth : 500;
            
            // 重置位置实现无缝循环
            if (state.marqueeX < -state.marqueeWidth) {
                state.marqueeX = containerWidth;
            }
            
            // 应用transform（兼容旧版浏览器）
            var transform = 'translateX(' + state.marqueeX + 'px)';
            elements.marqueeText.style.webkitTransform = transform;
            elements.marqueeText.style.mozTransform = transform;
            elements.marqueeText.style.msTransform = transform;
            elements.marqueeText.style.oTransform = transform;
            elements.marqueeText.style.transform = transform;
        }, frameTime);
        
        // 时间更新
        state.timeTimer = setInterval(function() {
            updateTime();
        }, 1000);
    }

    /**
     * 停止所有动画
     */
    function stopAnimations() {
        if (state.scrollTimer) {
            clearInterval(state.scrollTimer);
            state.scrollTimer = null;
        }
        if (state.marqueeTimer) {
            clearInterval(state.marqueeTimer);
            state.marqueeTimer = null;
        }
        if (state.timeTimer) {
            clearInterval(state.timeTimer);
            state.timeTimer = null;
        }
    }

    /**
     * 更新时间显示
     */
    function updateTime() {
        if (!elements.datetime) return;
        
        var now = new Date();
        var days = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'];
        
        var year = now.getFullYear();
        var month = padZero(now.getMonth() + 1);
        var date = padZero(now.getDate());
        var day = days[now.getDay()];
        var hours = padZero(now.getHours());
        var minutes = padZero(now.getMinutes());
        var seconds = padZero(now.getSeconds());
        
        elements.datetime.textContent = year + '-' + month + '-' + date + ' ' + day + ' ' + hours + ':' + minutes + ':' + seconds;
    }

    /**
     * 数字补零
     */
    function padZero(num) {
        return (num < 10 ? '0' : '') + num;
    }

    /**
     * 启动刷新定时器
     */
    function startRefreshTimer() {
        state.refreshTimer = setInterval(function() {
            console.log('🔄 定时刷新金价数据...');
            fetchPriceData();
        }, CONFIG.refreshInterval);
    }

    /**
     * 清理资源
     */
    function cleanup() {
        stopAnimations();
        if (state.refreshTimer) {
            clearInterval(state.refreshTimer);
            state.refreshTimer = null;
        }
    }

    // 页面加载完成后初始化
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // 页面卸载时清理
    window.addEventListener('beforeunload', cleanup);

})();

