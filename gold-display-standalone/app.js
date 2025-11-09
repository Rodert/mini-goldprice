/**
 * 金价展示页面 - 原生 JavaScript 实现
 * 兼容 32位低版本 Fully Kiosk Browser (Android 4.4+, Chrome 30+)
 * 使用 ES5 语法，不使用 async/await、箭头函数等新特性
 */

(function() {
    'use strict';

    // 配置
    var CONFIG = {
        apiBaseUrl: '/api/v1',
        scrollSpeed: 30,      // 产品图片滚动速度（像素/秒）
        marqueeSpeed: 100,    // 欢迎语滚动速度（像素/秒）
        refreshInterval: 5 * 60 * 1000,  // 价格刷新间隔（5分钟）
        productImages: [
            'assets/img.jpg',
            'assets/img.jpg',
            'assets/img.jpg',
            'assets/img.jpg',
            'assets/img.jpg'
        ]
    };

    // 应用状态
    var app = {
        scrollY: 0,
        marqueeX: 0,
        marqueeWidth: 0,
        scrollTimer: null,
        marqueeTimer: null,
        timeTimer: null,
        refreshTimer: null,
        priceList: []
    };

    // DOM 元素引用
    var elements = {};

    /**
     * 初始化应用
     */
    function init() {
        // 获取 DOM 元素
        elements.videoPlayer = document.getElementById('videoPlayer');
        elements.videoPlaceholder = document.getElementById('videoPlaceholder');
        elements.scrollContainer = document.getElementById('scrollContainer');
        elements.priceTable = document.getElementById('priceTable');
        elements.priceTableBody = document.getElementById('priceTableBody');
        elements.loadingMsg = document.getElementById('loadingMsg');
        elements.datetime = document.getElementById('datetime');
        elements.marqueeText = document.getElementById('marqueeText');

        // 初始化各个模块
        initVideo();
        initProductImages();
        initMarquee();
        startAnimations();
        updateTime();
        fetchPriceData();

        // 定期刷新价格数据
        app.refreshTimer = setInterval(function() {
            fetchPriceData();
        }, CONFIG.refreshInterval);

        // 定期更新时间
        app.timeTimer = setInterval(function() {
            updateTime();
        }, 1000);
    }

    /**
     * 初始化视频
     */
    function initVideo() {
        if (!elements.videoPlayer) return;

        // 检查视频是否可以播放
        elements.videoPlayer.addEventListener('error', function() {
            elements.videoPlayer.style.display = 'none';
            if (elements.videoPlaceholder) {
                elements.videoPlaceholder.style.display = 'flex';
            }
        });

        // 视频加载完成
        elements.videoPlayer.addEventListener('loadedmetadata', function() {
            console.log('视频加载完成');
        });
    }

    /**
     * 初始化产品图片
     */
    function initProductImages() {
        if (!elements.scrollContainer) return;

        // 创建双倍图片数组用于无缝滚动
        var doubleImages = CONFIG.productImages.concat(CONFIG.productImages);
        var html = '';

        for (var i = 0; i < doubleImages.length; i++) {
            html += '<div class="product-image">';
            html += '<img src="' + doubleImages[i] + '" alt="产品图片" onerror="this.src=\'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iIzMzMzMzMyIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBmb250LXNpemU9IjE0IiBmaWxsPSIjOTk5IiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBkeT0iLjNlbSI+5Zu+54mH5Yqg6L29PC90ZXh0Pjwvc3ZnPg==\'">';
            html += '</div>';
        }

        elements.scrollContainer.innerHTML = html;
    }

    /**
     * 初始化跑马灯
     */
    function initMarquee() {
        if (!elements.marqueeText) return;

        // 等待 DOM 渲染完成
        setTimeout(function() {
            if (elements.marqueeText) {
                app.marqueeWidth = elements.marqueeText.offsetWidth || 500;
                var marqueeContent = elements.marqueeText.parentElement;
                if (marqueeContent) {
                    app.marqueeX = marqueeContent.offsetWidth || 500;
                }
            }
        }, 100);
    }

    /**
     * 开始所有动画
     */
    function startAnimations() {
        // 产品图片垂直滚动
        app.scrollTimer = setInterval(function() {
            app.scrollY -= CONFIG.scrollSpeed / 60;
            var singleHeight = CONFIG.productImages.length * 220; // 每张图200px + 20px padding
            if (Math.abs(app.scrollY) >= singleHeight) {
                app.scrollY = 0;
            }
            if (elements.scrollContainer) {
                elements.scrollContainer.style.transform = 'translateY(' + app.scrollY + 'px)';
                elements.scrollContainer.style.webkitTransform = 'translateY(' + app.scrollY + 'px)';
            }
        }, 1000 / 60);

        // 跑马灯横向滚动
        app.marqueeTimer = setInterval(function() {
            app.marqueeX -= CONFIG.marqueeSpeed / 60;
            var marqueeContent = elements.marqueeText ? elements.marqueeText.parentElement : null;
            var containerWidth = marqueeContent ? marqueeContent.offsetWidth : 500;
            if (app.marqueeX < -app.marqueeWidth) {
                app.marqueeX = containerWidth;
            }
            if (elements.marqueeText) {
                elements.marqueeText.style.transform = 'translateX(' + app.marqueeX + 'px)';
                elements.marqueeText.style.webkitTransform = 'translateX(' + app.marqueeX + 'px)';
            }
        }, 1000 / 60);
    }

    /**
     * 停止所有动画
     */
    function stopAnimations() {
        if (app.scrollTimer) {
            clearInterval(app.scrollTimer);
            app.scrollTimer = null;
        }
        if (app.marqueeTimer) {
            clearInterval(app.marqueeTimer);
            app.marqueeTimer = null;
        }
        if (app.timeTimer) {
            clearInterval(app.timeTimer);
            app.timeTimer = null;
        }
        if (app.refreshTimer) {
            clearInterval(app.refreshTimer);
            app.refreshTimer = null;
        }
    }

    /**
     * 更新时间显示
     */
    function updateTime() {
        var now = new Date();
        var days = ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'];
        var year = now.getFullYear();
        var month = padZero(now.getMonth() + 1);
        var date = padZero(now.getDate());
        var day = days[now.getDay()];
        var hours = padZero(now.getHours());
        var minutes = padZero(now.getMinutes());
        var seconds = padZero(now.getSeconds());

        if (elements.datetime) {
            elements.datetime.textContent = year + '-' + month + '-' + date + ' ' + day + ' ' + hours + ':' + minutes + ':' + seconds;
        }
    }

    /**
     * 补零函数
     */
    function padZero(num) {
        return (num < 10 ? '0' : '') + num;
    }

    /**
     * 获取金价数据
     */
    function fetchPriceData() {
        var url = CONFIG.apiBaseUrl + '/prices?page=1&page_size=100';
        var xhr = new XMLHttpRequest();

        xhr.open('GET', url, true);
        xhr.setRequestHeader('Content-Type', 'application/json');

        xhr.onreadystatechange = function() {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    try {
                        var response = JSON.parse(xhr.responseText);
                        handlePriceResponse(response);
                    } catch (e) {
                        console.error('解析响应失败:', e);
                        useMockData();
                    }
                } else {
                    console.warn('API 请求失败，状态码:', xhr.status);
                    useMockData();
                }
            }
        };

        xhr.onerror = function() {
            console.error('网络错误');
            useMockData();
        };

        xhr.send();
    }

    /**
     * 处理价格响应
     */
    function handlePriceResponse(response) {
        var dataList = null;

        // 尝试不同的数据结构
        if (response && response.code === 200) {
            if (response.data && response.data.list) {
                dataList = response.data.list;
            } else if (response.data && Array.isArray(response.data)) {
                dataList = response.data;
            } else if (Array.isArray(response)) {
                dataList = response;
            }
        }

        if (dataList && dataList.length > 0) {
            app.priceList = dataList.slice(0, 5).map(function(item) {
                return {
                    id: item.id,
                    name: item.name,
                    price: item.sell_price || item.price || 0,
                    fee: item.fee || 10
                };
            });
            renderPriceTable();
        } else {
            console.warn('数据列表为空，使用模拟数据');
            useMockData();
        }
    }

    /**
     * 使用模拟数据
     */
    function useMockData() {
        app.priceList = [
            { id: 1, name: '足金(9999)', price: 488, fee: 10 },
            { id: 2, name: '足金(999)', price: 428, fee: 10 },
            { id: 3, name: 'Pt950', price: 388, fee: 10 },
            { id: 4, name: 'Pt990', price: 408, fee: 10 },
            { id: 5, name: 'PD950', price: 218, fee: 10 }
        ];
        renderPriceTable();
    }

    /**
     * 渲染价格表格
     */
    function renderPriceTable() {
        if (!elements.priceTableBody) return;

        var html = '';
        for (var i = 0; i < app.priceList.length; i++) {
            var item = app.priceList[i];
            html += '<tr>';
            html += '<td class="product-name">' + escapeHtml(item.name) + '</td>';
            html += '<td class="price">' + escapeHtml(item.price) + '</td>';
            html += '<td class="fee">' + escapeHtml(item.fee) + '</td>';
            html += '</tr>';
        }

        elements.priceTableBody.innerHTML = html;

        // 显示表格，隐藏加载消息
        if (elements.loadingMsg) {
            elements.loadingMsg.style.display = 'none';
        }
        if (elements.priceTable) {
            elements.priceTable.style.display = 'table';
        }
    }

    /**
     * HTML 转义
     */
    function escapeHtml(text) {
        var map = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        return String(text).replace(/[&<>"']/g, function(m) {
            return map[m];
        });
    }

    // 页面加载完成后初始化
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }

    // 页面卸载时清理
    window.addEventListener('beforeunload', function() {
        stopAnimations();
    });

})();

