package m3u8_test

import (
	"log"
	"testing"

	"github.com/AmbitiousJun/live-server/internal/service/m3u8"
)

const content = `#EXTM3U
#EXTINF:-1 tvg-name="CCTV1" tvg-id="CCTV1" tvg-logo="https://live.fanmingming.com/tv/CCTV1.png" group-title="央视频道", CCTV1
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010144/1.m3u8
#EXTINF:-1 tvg-name="CCTV2" tvg-id="CCTV2" tvg-logo="https://live.fanmingming.com/tv/CCTV2.png" group-title="央视频道", CCTV2
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000023315/1.m3u8
#EXTINF:-1 tvg-name="CCTV3" tvg-id="CCTV3" tvg-logo="https://live.fanmingming.com/tv/CCTV3.png" group-title="央视频道", CCTV3
http://[2409:8087:5e00:24::10]:6060/200000001898/460000089800010212/main.m3u8
#EXTINF:-1 tvg-name="CCTV4" tvg-id="CCTV4" tvg-logo="https://live.fanmingming.com/tv/CCTV4.png" group-title="央视频道", CCTV4
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000005000265004/1.m3u8
#EXTINF:-1 tvg-name="CCTV5" tvg-id="CCTV5" tvg-logo="https://live.fanmingming.com/tv/CCTV5.png" group-title="央视频道", CCTV5
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000004794/1.m3u8
#EXTINF:-1 tvg-name="CCTV5+" tvg-id="CCTV5+" tvg-logo="https://live.fanmingming.com/tv/CCTV5+.png" group-title="央视频道", CCTV5+
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000018504/1.m3u8
#EXTINF:-1 tvg-name="CCTV6" tvg-id="CCTV6" tvg-logo="https://live.fanmingming.com/tv/CCTV6.png" group-title="央视频道", CCTV6
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000016466/1.m3u8
#EXTINF:-1 tvg-name="CCTV7" tvg-id="CCTV7" tvg-logo="https://live.fanmingming.com/tv/CCTV7.png" group-title="央视频道", CCTV7
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010215/1.m3u8
#EXTINF:-1 tvg-name="CCTV8" tvg-id="CCTV8" tvg-logo="https://live.fanmingming.com/tv/CCTV8.png" group-title="央视频道", CCTV8
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000003736/1.m3u8
#EXTINF:-1 tvg-name="CCTV9" tvg-id="CCTV9" tvg-logo="https://live.fanmingming.com/tv/CCTV9.png" group-title="央视频道", CCTV9
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000028286/1.m3u8
#EXTINF:-1 tvg-name="CCTV10" tvg-id="CCTV10" tvg-logo="https://live.fanmingming.com/tv/CCTV10.png" group-title="央视频道", CCTV10
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000012827/1.m3u8
#EXTINF:-1 tvg-name="CCTV11" tvg-id="CCTV11" tvg-logo="https://live.fanmingming.com/tv/CCTV11.png" group-title="央视频道", CCTV11
http://[2409:8087:74d9:21::6]:80/270000001322/69900158041111100000002177/index.m3u8
#EXTINF:-1 tvg-name="CCTV12" tvg-id="CCTV12" tvg-logo="https://live.fanmingming.com/tv/CCTV12.png" group-title="央视频道", CCTV12
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031107/1.m3u8
#EXTINF:-1 tvg-name="CCTV13" tvg-id="CCTV13" tvg-logo="https://live.fanmingming.com/tv/CCTV13.png" group-title="央视频道", CCTV13
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000021303/1.m3u8
#EXTINF:-1 tvg-name="CCTV14" tvg-id="CCTV14" tvg-logo="https://live.fanmingming.com/tv/CCTV14.png" group-title="央视频道", CCTV14
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000010000023358/1.m3u8
#EXTINF:-1 tvg-name="CCTV15" tvg-id="CCTV15" tvg-logo="https://live.fanmingming.com/tv/CCTV15.png" group-title="央视频道", CCTV15
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031109/1.m3u8
#EXTINF:-1 tvg-name="CCTV16" tvg-id="CCTV16" tvg-logo="https://live.fanmingming.com/tv/CCTV16.png" group-title="央视频道", CCTV16
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010142/index.m3u8
#EXTINF:-1 tvg-name="CCTV17" tvg-id="CCTV17" tvg-logo="https://live.fanmingming.com/tv/CCTV17.png" group-title="央视频道", CCTV17
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010210/1.m3u8
#EXTINF:-1 tvg-name="CCTV中视购物" tvg-id="CCTV中视购物" tvg-logo="" group-title="央视频道", CCTV中视购物
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010134/1.m3u8
#EXTINF:-1 tvg-name="安徽卫视" tvg-id="安徽卫视" tvg-logo="https://live.fanmingming.com/tv/安徽卫视.png" group-title="卫视频道", 安徽卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000023002/1.m3u8
#EXTINF:-1 tvg-name="北京卫视" tvg-id="北京卫视" tvg-logo="https://live.fanmingming.com/tv/北京卫视.png" group-title="卫视频道", 北京卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000021288/1.m3u8
#EXTINF:-1 tvg-name="兵团卫视" tvg-id="兵团卫视" tvg-logo="https://live.fanmingming.com/tv/兵团卫视.png" group-title="卫视频道", 兵团卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000005000266005/1.m3u8
#EXTINF:-1 tvg-name="东方卫视" tvg-id="东方卫视" tvg-logo="https://live.fanmingming.com/tv/东方卫视.png" group-title="卫视频道", 东方卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000014098/1.m3u8
#EXTINF:-1 tvg-name="东南卫视" tvg-id="东南卫视" tvg-logo="https://live.fanmingming.com/tv/东南卫视.png" group-title="卫视频道", 东南卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000010584/1.m3u8
#EXTINF:-1 tvg-name="甘肃卫视" tvg-id="甘肃卫视" tvg-logo="" group-title="卫视频道", 甘肃卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031121/1.m3u8
#EXTINF:-1 tvg-name="广东卫视" tvg-id="广东卫视" tvg-logo="https://live.fanmingming.com/tv/广东卫视.png" group-title="卫视频道", 广东卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000014176/1.m3u8
#EXTINF:-1 tvg-name="广西卫视" tvg-id="广西卫视" tvg-logo="https://live.fanmingming.com/tv/广西卫视.png" group-title="卫视频道", 广西卫视
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010099/1.m3u8
#EXTINF:-1 tvg-name="贵州卫视" tvg-id="贵州卫视" tvg-logo="https://live.fanmingming.com/tv/贵州卫视.png" group-title="卫视频道", 贵州卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000003169/1.m3u8
#EXTINF:-1 tvg-name="海南卫视" tvg-id="海南卫视" tvg-logo="https://live.fanmingming.com/tv/海南卫视.png" group-title="卫视频道", 海南卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000006211/index.m3u8
#EXTINF:-1 tvg-name="河北卫视" tvg-id="河北卫视" tvg-logo="https://live.fanmingming.com/tv/河北卫视.png" group-title="卫视频道", 河北卫视
http://[2409:8087:74d9:21::6]:80/000000001000PLTV/88888888/224/3221226442/index.m3u8
#EXTINF:-1 tvg-name="河南卫视" tvg-id="河南卫视" tvg-logo="https://live.fanmingming.com/tv/河南卫视.png" group-title="卫视频道", 河南卫视
http://[2409:8087:74d9:21::6]:80/000000001000PLTV/88888888/224/3221226614/index.m3u8
#EXTINF:-1 tvg-name="湖北卫视" tvg-id="湖北卫视" tvg-logo="https://live.fanmingming.com/tv/湖北卫视.png" group-title="卫视频道", 湖北卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000014954/1.m3u8
#EXTINF:-1 tvg-name="湖南卫视" tvg-id="湖南卫视" tvg-logo="https://live.fanmingming.com/tv/湖南卫视.png" group-title="卫视频道", 湖南卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000006692/1.m3u8
#EXTINF:-1 tvg-name="吉林卫视" tvg-id="吉林卫视" tvg-logo="https://live.fanmingming.com/tv/吉林卫视.png" group-title="卫视频道", 吉林卫视
http://[2409:8087:74d9:21::6]:80/000000001000PLTV/88888888/224/3221226440/index.m3u8
#EXTINF:-1 tvg-name="江苏卫视" tvg-id="江苏卫视" tvg-logo="https://live.fanmingming.com/tv/江苏卫视.png" group-title="卫视频道", 江苏卫视
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010059/index.m3u8
#EXTINF:-1 tvg-name="江西卫视" tvg-id="江西卫视" tvg-logo="https://live.fanmingming.com/tv/江西卫视.png" group-title="卫视频道", 江西卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000006000268001/1.m3u8
#EXTINF:-1 tvg-name="辽宁卫视" tvg-id="辽宁卫视" tvg-logo="" group-title="卫视频道", 辽宁卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000011671/1.m3u8
#EXTINF:-1 tvg-name="宁夏卫视" tvg-id="宁夏卫视" tvg-logo="https://live.fanmingming.com/tv/宁夏卫视.png" group-title="卫视频道", 宁夏卫视
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010110/index.m3u8
#EXTINF:-1 tvg-name="青海卫视" tvg-id="青海卫视" tvg-logo="" group-title="卫视频道", 青海卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000013359/1.m3u8
#EXTINF:-1 tvg-name="三沙卫视" tvg-id="三沙卫视" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/三沙卫视.png" group-title="卫视频道", 三沙卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/4600001000000000117/1.m3u8
#EXTINF:-1 tvg-name="厦门卫视" tvg-id="厦门卫视" tvg-logo="https://live.fanmingming.com/tv/厦门卫视.png" group-title="卫视频道", 厦门卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000005000266006/1.m3u8
#EXTINF:-1 tvg-name="山东卫视" tvg-id="山东卫视" tvg-logo="https://live.fanmingming.com/tv/山东卫视.png" group-title="卫视频道", 山东卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000016568/1.m3u8
#EXTINF:-1 tvg-name="山西卫视" tvg-id="山西卫视" tvg-logo="https://live.fanmingming.com/tv/山西卫视.png" group-title="卫视频道", 山西卫视
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010114/1.m3u8
#EXTINF:-1 tvg-name="陕西卫视" tvg-id="陕西卫视" tvg-logo="" group-title="卫视频道", 陕西卫视
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010113/1.m3u8
#EXTINF:-1 tvg-name="深圳卫视" tvg-id="深圳卫视" tvg-logo="https://live.fanmingming.com/tv/深圳卫视.png" group-title="卫视频道", 深圳卫视
http://[2409:8087:74d9:21::6]:80/000000001000PLTV/88888888/224/3221226199/index.m3u8
#EXTINF:-1 tvg-name="四川卫视" tvg-id="四川卫视" tvg-logo="https://live.fanmingming.com/tv/四川卫视.png" group-title="卫视频道", 四川卫视
http://[2409:8087:74d9:21::6]:80/000000001000PLTV/88888888/224/3221226454/index.m3u8
#EXTINF:-1 tvg-name="天津卫视" tvg-id="天津卫视" tvg-logo="https://live.fanmingming.com/tv/天津卫视.png" group-title="卫视频道", 天津卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000000831/1.m3u8
#EXTINF:-1 tvg-name="西藏卫视" tvg-id="西藏卫视" tvg-logo="https://live.fanmingming.com/tv/西藏卫视.png" group-title="卫视频道", 西藏卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/6603041244077933770/1.m3u8
#EXTINF:-1 tvg-name="新疆卫视" tvg-id="新疆卫视" tvg-logo="https://live.fanmingming.com/tv/新疆卫视.png" group-title="卫视频道", 新疆卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000029441/1.m3u8
#EXTINF:-1 tvg-name="延边卫视" tvg-id="延边卫视" tvg-logo="https://live.fanmingming.com/tv/延边卫视.png" group-title="卫视频道", 延边卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000005000266008/1.m3u8
#EXTINF:-1 tvg-name="云南卫视" tvg-id="云南卫视" tvg-logo="https://live.fanmingming.com/tv/云南卫视.png" group-title="卫视频道", 云南卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000024694/1.m3u8
#EXTINF:-1 tvg-name="浙江卫视" tvg-id="浙江卫视" tvg-logo="https://live.fanmingming.com/tv/浙江卫视.png" group-title="卫视频道", 浙江卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000001000009806/1.m3u8
#EXTINF:-1 tvg-name="重庆卫视" tvg-id="重庆卫视" tvg-logo="" group-title="卫视频道", 重庆卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000018937/1.m3u8
#EXTINF:-1 tvg-name="黑龙江卫视" tvg-id="黑龙江卫视" tvg-logo="https://live.fanmingming.com/tv/黑龙江卫视.png" group-title="卫视频道", 黑龙江卫视
http://[2409:8087:74d9:21::6]:80/270000001322/69900158041111100000002142/index.m3u8
#EXTINF:-1 tvg-name="内蒙古卫视" tvg-id="内蒙古卫视" tvg-logo="https://live.fanmingming.com/tv/内蒙古卫视.png" group-title="卫视频道", 内蒙古卫视
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010109/1.m3u8
#EXTINF:-1 tvg-name="东方卫视" tvg-id="东方卫视" tvg-logo="https://live.fanmingming.com/tv/东方卫视.png" group-title="卫视频道", 东方卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000004000014098/1.m3u8
#EXTINF:-1 tvg-name="山东教育" tvg-id="山东教育" tvg-logo="https://gitee.com/Black_crow/epglogo/raw/master/shandong/Shandong9.png" group-title="卫视频道", 山东教育
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010112/1.m3u8
#EXTINF:-1 tvg-name="CHC动作电影" tvg-id="CHC动作电影" tvg-logo="https://live.fanmingming.com/tv/CHC动作电影.png" group-title="电影频道", CHC动作电影
http://[2409:8087:5e00:24::1e]:6060/000000001000/8103864434730665389/1.m3u8
#EXTINF:-1 tvg-name="CHC影迷电影" tvg-id="CHC影迷电影" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/CHC影迷电影.png" group-title="电影频道", CHC影迷电影
http://[2001:250:5800:1005::155]:80/liverespath/bca6167afdef7fc773405c38e695b342c1d1eca0/index.m3u8
#EXTINF:-1 tvg-name="黑莓电影" tvg-id="黑莓电影" tvg-logo="https://live.fanmingming.com/tv/黑莓电影.png" group-title="电影频道", 黑莓电影
http://[2409:8087:74d9:21::6]:80/270000001128/9900000095/index.m3u8
#EXTINF:-1 tvg-name="东方影视" tvg-id="东方影视" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/东方影视-light.png" group-title="电影频道", 东方影视
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000002000000013/index.m3u8
#EXTINF:-1 tvg-name="黑龙江影视" tvg-id="黑龙江影视" tvg-logo="" group-title="电影频道", 黑龙江影视
http://[2409:8087:1a01:df::4051]:80/TVOD/88888888/224/3221225973/main.m3u8
#EXTINF:-1 tvg-name="超级电影" tvg-id="超级电影" tvg-logo="" group-title="电影频道", 超级电影
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000003000012426/1.m3u8
#EXTINF:-1 tvg-name="CHC家庭影院" tvg-id="CHC家庭影院" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/CHC家庭影院.png" group-title="数字频道", CHC家庭影院
http://[2001:250:5800:1005::155]:80/liverespath/449d51febadb152094085d373b9af94a6ac5f1dd/index.m3u8
#EXTINF:-1 tvg-name="嘉佳卡通" tvg-id="嘉佳卡通" tvg-logo="https://live.fanmingming.com/tv/嘉佳卡通.png" group-title="数字频道", 嘉佳卡通
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000025964/1.m3u8
#EXTINF:-1 tvg-name="乐游频道" tvg-id="乐游频道" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/乐游频道.png" group-title="数字频道", 乐游频道
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031112/1.m3u8
#EXTINF:-1 tvg-name="SITV动漫秀场" tvg-id="动漫秀场" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/SITV动漫秀场-light.png" group-title="数字频道", SITV动漫秀场
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031113/1.m3u8
#EXTINF:-1 tvg-name="SITV都市剧场" tvg-id="都市剧场" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/SITV都市剧场-light.png" group-title="数字频道", SITV都市剧场
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031111/1.m3u8
#EXTINF:-1 tvg-name="SITV欢笑剧场" tvg-id="欢笑剧场" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/SITV欢笑剧场-light.png" group-title="数字频道", SITV欢笑剧场
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000002000009455/1.m3u8
#EXTINF:-1 tvg-name="SITV金色学堂" tvg-id="金色学堂" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/SITV金色学堂.png" group-title="数字频道", SITV金色学堂
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000010000026105/1.m3u8
#EXTINF:-1 tvg-name="游戏风云" tvg-id="游戏风云" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/游戏风云.png" group-title="数字频道", 游戏风云
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031114/1.m3u8
#EXTINF:-1 tvg-name="哒啵电竞" tvg-id="哒啵电竞" tvg-logo="https://live.fanmingming.com/tv/哒啵电竞.png" group-title="数字频道", 哒啵电竞
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000003000000066/index.m3u8
#EXTINF:-1 tvg-name="劲爆体育" tvg-id="劲爆体育" tvg-logo="https://gitee.com/Black_crow/epglogo/raw/master/shuzi/jbty.png" group-title="数字频道", 劲爆体育
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000002000000008/index.m3u8
#EXTINF:-1 tvg-name="五星体育" tvg-id="五星体育" tvg-logo="https://live.fanmingming.com/tv/五星体育.png" group-title="数字频道", 五星体育
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000002000000007/index.m3u8
#EXTINF:-1 tvg-name="五星体育" tvg-id="五星体育" tvg-logo="https://live.fanmingming.com/tv/五星体育.png" group-title="数字频道", 五星体育
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000002000000007/index.m3u8
#EXTINF:-1 tvg-name="快乐垂钓" tvg-id="快乐垂钓" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/快乐垂钓.png" group-title="数字频道", 快乐垂钓
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031206/1.m3u8
#EXTINF:-1 tvg-name="求索动物" tvg-id="求索动物" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/求索动物-light.png" group-title="数字频道", 求索动物
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000004000000009/index.m3u8
#EXTINF:-1 tvg-name="求索纪录" tvg-id="求索纪录" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/求索纪录-light.png" group-title="数字频道", 求索纪录
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000004000000010/index.m3u8
#EXTINF:-1 tvg-name="求索科学" tvg-id="求索科学" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/求索科学-light.png" group-title="数字频道", 求索科学
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000004000000011/index.m3u8
#EXTINF:-1 tvg-name="求索生活" tvg-id="171" tvg-logo="https://gitee.com/Black_crow/epglogo/raw/master/shuzi/Qiusuo4.png" group-title="数字频道", 求索生活
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000004000000008/index.m3u8
#EXTINF:-1 tvg-name="法治天地" tvg-id="法治天地" tvg-logo="https://live.fanmingming.com/tv/法治天地.png" group-title="数字频道", 法治天地
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000002000000014/index.m3u8
#EXTINF:-1 tvg-name="生活时尚" tvg-id="生活时尚" tvg-logo="https://live.fanmingming.com/tv/生活时尚.png" group-title="数字频道", 生活时尚
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000002000000006/index.m3u8
#EXTINF:-1 tvg-name="茶友频道" tvg-id="茶友频道" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/茶友频道.png" group-title="数字频道", 茶友频道
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031209/1.m3u8
#EXTINF:-1 tvg-name="纯享4K" tvg-id="纯享4K" tvg-logo="" group-title="数字频道", 纯享4K
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000004000011651/1.m3u8
#EXTINF:-1 tvg-name="广东珠江" tvg-id="广东珠江" tvg-logo="https://live.fanmingming.com/tv/广东珠江.png" group-title="地方频道", 广东珠江
http://[2409:8087:5e00:24::1e]:6060/200000001898/460000089800010091/1.m3u8
#EXTINF:-1 tvg-name="大湾区卫视" tvg-id="大湾区卫视" tvg-logo="https://live.fanmingming.com/tv/大湾区卫视.png" group-title="地方频道", 大湾区卫视
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000002000011619/1.m3u8
#EXTINF:-1 tvg-name="金鹰纪实" tvg-id="金鹰纪实" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/金鹰纪实.png" group-title="地方频道", 金鹰纪实
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031203/1.m3u8
#EXTINF:-1 tvg-name="上海新闻综合" tvg-id="上海新闻" tvg-logo="https://live.fanmingming.com/tv/上视新闻.png" group-title="地方频道", 上海新闻综合
http://[2409:8087:5e08:24::11]:6610/000000001000/2000000002000000005/index.m3u8
#EXTINF:-1 tvg-name="上海外语" tvg-id="上海外语" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/上海外语-light.png" group-title="地方频道", 上海外语
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000002000000001/index.m3u8
#EXTINF:-1 tvg-name="上海新闻" tvg-id="上海新闻" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/上海新闻-light.png" group-title="地方频道", 上海新闻
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000002000000005/index.m3u8
#EXTINF:-1 tvg-name="第一财经" tvg-id="第一财经" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/第一财经-light.png" group-title="地方频道", 第一财经
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000002000000004/index.m3u8
#EXTINF:-1 tvg-name="东方财经" tvg-id="东方财经" tvg-logo="" group-title="地方频道", 东方财经
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000007000010003/1.m3u8
#EXTINF:-1 tvg-name="七彩戏剧" tvg-id="生活时尚" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/七彩戏剧-light.png" group-title="地方频道", 七彩戏剧
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000002000000010/index.m3u8
#EXTINF:-1 tvg-name="上海都市" tvg-id="上海都市" tvg-logo="https://mirror.ghproxy.com/https://raw.githubusercontent.com/drangjchen/IPTV/main/Logo/上海都市-light.png" group-title="地方频道", 上海都市
http://[2409:8087:5e01:24::16]:6610/000000001000/2000000002000000012/index.m3u8
#EXTINF:-1 tvg-name="海南新闻" tvg-id="海南新闻" tvg-logo="https://gitee.com/Black_crow/epglogo/raw/master/hainan/Hainan2.png" group-title="地方频道", 海南新闻
http://[2409:8087:5e00:24::1e]:6060/000000001000/4600001000000000111/1.m3u8
#EXTINF:-1 tvg-name="海南自贸" tvg-id="海南自贸" tvg-logo="https://gitee.com/Black_crow/epglogo/raw/master/hainan/Hainan6.png" group-title="地方频道", 海南自贸
http://[2409:8087:5e00:24::1e]:6060/000000001000/4600001000000000116/1.m3u8
#EXTINF:-1 tvg-name="海南文旅" tvg-id="海南文旅" tvg-logo="https://live.fanmingming.com/tv/海南文旅.png" group-title="地方频道", 海南文旅
http://[2409:8087:5e00:24::1e]:6060/000000001000/4600001000000000113/1.m3u8
#EXTINF:-1 tvg-name="海南公共" tvg-id="海南公共" tvg-logo="https://live.fanmingming.com/tv/海南公共.png" group-title="地方频道", 海南公共
http://[2409:8087:5e00:24::1e]:6060/000000001000/460000100000000057/1.m3u8
#EXTINF:-1 tvg-name="黑龙江文体" tvg-id="黑龙江文体" tvg-logo="" group-title="地方频道", 黑龙江文体
http://[2409:8087:1a01:df::4025]:80/TVOD/88888888/224/3221225965/main.m3u8
#EXTINF:-1 tvg-name="黑龙江都市" tvg-id="黑龙江都市" tvg-logo="" group-title="地方频道", 黑龙江都市
http://[2409:8087:1a01:df::4055]:80/TVOD/88888888/224/3221225969/main.m3u8
#EXTINF:-1 tvg-name="鹤岗新闻综合" tvg-id="鹤岗新闻综合" tvg-logo="" group-title="地方频道", 鹤岗新闻综合
http://[2409:8087:1a01:df::4059]:80/TVOD/88888888/224/3221226073/main.m3u8
#EXTINF:-1 tvg-name="黑龙江农业科教" tvg-id="HLJ05" tvg-logo="https://gitee.com/suxuang/TVlogo/raw/main/img/Heilongjiang5.png" group-title="地方频道", 黑龙江农业科教
http://[2409:8087:1a01:df::4077]:80/PLTV/88888888/224/3221225994/index.m3u8
#EXTINF:-1 tvg-name="黑龙江新闻" tvg-id="HLJ01" tvg-logo="https://gitee.com/suxuang/TVlogo/raw/main/img/Heilongjiang1.png" group-title="地方频道", 黑龙江新闻
http://[2409:8087:1a01:df::4077]:80/PLTV/88888888/224/3221225967/index.m3u8
#EXTINF:-1 tvg-name="超级综艺" tvg-id="超级综艺" tvg-logo="" group-title="综艺频道", 超级综艺
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000006000268002/1.m3u8
#EXTINF:-1 tvg-name="金牌综艺" tvg-id="金牌综艺" tvg-logo="" group-title="综艺频道", 金牌综艺
http://[2409:8087:5e00:24::1e]:6060/000000001000/6399725674632152632/1.m3u8
#EXTINF:-1 tvg-name="精品体育" tvg-id="精品体育" tvg-logo="https://live.fanmingming.com/tv/NEWTV精品体育.png" group-title="体育频道", 精品体育
http://[2409:8087:74d9:21::6]:80/270000001128/9900000102/index.m3u8
#EXTINF:-1 tvg-name="睛彩竞技" tvg-id="睛彩竞技" tvg-logo="https://ghproxy.net/https://raw.githubusercontent.com/kimwang1978/collect-tv-txt/main/assets/logo/睛彩竞技.png" group-title="体育频道", 睛彩竞技
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000006000270005/1.m3u8
#EXTINF:-1 tvg-name="睛彩篮球" tvg-id="睛彩篮球" tvg-logo="https://ghproxy.net/https://raw.githubusercontent.com/kimwang1978/collect-tv-txt/main/assets/logo/睛彩篮球.png" group-title="体育频道", 睛彩篮球
http://[2409:8087:5e00:24::1e]:6060/000000001000/1000000006000270006/1.m3u8
#EXTINF:-1 tvg-name="魅力足球" tvg-id="魅力足球" tvg-logo="" group-title="体育频道", 魅力足球
http://[2409:8087:5e00:24::1e]:6060/000000001000/5000000011000031207/1.m3u8
#EXTINF:-1 tvg-name="2024-10-11 16:23:05" tvg-id="更新时间" tvg-logo="" group-title="更新时间", 2024-10-11 16:23:05
https://vd2.bdstatic.com/mda-qiddtcxr5ktdh0uv/720p/h264/1726307097251597100/mda-qiddtcxr5ktdh0uv.mp4?v_from_s=bdapp-resbox-hna`

func TestReadContent(t *testing.T) {
	infos, err := m3u8.ReadContent(content)
	if err != nil {
		t.Fatal(err)
		return
	}
	log.Println(infos)
}
