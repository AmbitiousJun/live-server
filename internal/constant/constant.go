package constant

const (
	Version    = "v1.21.2"                                     // 版本号
	RepoAddr   = "https://github.com/AmbitiousJun/live-server" // 仓库地址
	HeadersSeg = "[[[:]]]"                                     // 分割请求头的分隔符
)

const (
	Dir_DataRoot = "data" // 数据存放根目录
	// Dir_DataRoot = "/Users/ambitious/Desktop/code/go/live-server/data" // 数据存放根目录, 开发环境
)

const (

	// HelpDocHtmlTemplate 帮助文档的 html 模板
	HelpDocHtmlTemplate = `PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KICA8aGVhZD4KICAgIDxtZXRhIGNoYXJzZXQ9IlVURi04IiAvPgogICAgPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xLjAiIC8+CiAgICA8dGl0bGU+bGl2ZS1zZXJ2ZXIg5biu5Yqp5paH5qGjPC90aXRsZT4KICAgIDxzdHlsZT4KICAgICAgLmNvbnRhaW5lciB7CiAgICAgICAgb3ZlcmZsb3ctd3JhcDogYnJlYWstd29yZDsKICAgICAgICB3b3JkLWJyZWFrOiBicmVhay1hbGw7CiAgICAgIH0KICAgICAgLmNvbnRhaW5lciBhOmxpbmssCiAgICAgIC5jb250YWluZXIgYTp2aXNpdGVkIHsKICAgICAgICBjb2xvcjogIzFhNzNlODsKICAgICAgICB0ZXh0LWRlY29yYXRpb246IG5vbmU7CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9oZWFkPgogIDxib2R5PgogICAgPGRpdiBjbGFzcz0iY29udGFpbmVyIj4ke2RvY0NvbnRlbnR9PC9kaXY+CgogICAgPHNjcmlwdD4KICAgICAgd2luZG93Lm9ubG9hZCA9ICgpID0+IHsKICAgICAgICBjb25zdCBjb250YWluZXIgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIuY29udGFpbmVyIik7CiAgICAgICAgY29udGFpbmVyLmlubmVySFRNTCA9IGNvbnRhaW5lci5pbm5lckhUTUwucmVwbGFjZSgKICAgICAgICAgIC9cJFx7Y2xpZW50T3JpZ2luXH0vZywKICAgICAgICAgIGxvY2F0aW9uLm9yaWdpbgogICAgICAgICk7CiAgICAgIH07CiAgICA8L3NjcmlwdD4KICA8L2JvZHk+CjwvaHRtbD4K`

	// FengAuthHtml 凤凰秀授权地址
	FengAuthHtml = `PCFET0NUWVBFIGh0bWw+CjxodG1sPgogIDxoZWFkPgogICAgPG1ldGEgY2hhcnNldD0iVVRGLTgiIC8+CiAgICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBpbml0aWFsLXNjYWxlPTEuMCIgLz4KICAgIDx0aXRsZT7lh6Tlh7Dnp4DnmbvlvZXmjojmnYM8L3RpdGxlPgogICAgPCEtLSDlvJXlhaUgQm9vdHN0cmFwIENTUyAtLT4KICAgIDxsaW5rCiAgICAgIGhyZWY9Imh0dHBzOi8vZmFzdGx5LmpzZGVsaXZyLm5ldC9ucG0vYm9vdHN0cmFwQDUuMy4wL2Rpc3QvY3NzL2Jvb3RzdHJhcC5taW4uY3NzIgogICAgICByZWw9InN0eWxlc2hlZXQiCiAgICAvPgogICAgPHN0eWxlPgogICAgICAucmVzdWx0LXdyYXAgewogICAgICAgIG1hcmdpbi10b3A6IDQwcHg7CiAgICAgICAgcGFkZGluZzogMjBweDsKICAgICAgfQogICAgICAucmVzdWx0LXdyYXAgLm1lc3NhZ2UgewogICAgICAgIHBhZGRpbmc6IDEwcHg7CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9oZWFkPgogIDxib2R5PgogICAgPGRpdiBjbGFzcz0iY29udGFpbmVyIG10LTQiPgogICAgICA8IS0tIOaJi+acuuWPt+i+k+WFpeahhiAtLT4KICAgICAgPGRpdiBjbGFzcz0ibWItMyByb3ciPgogICAgICAgIDxsYWJlbCBmb3I9InBob25lIiBjbGFzcz0iY29sLXNtLTIgY29sLWZvcm0tbGFiZWwiPuaJi+acuuWPtzwvbGFiZWw+CiAgICAgICAgPGRpdiBjbGFzcz0iY29sLXNtLTEwIj4KICAgICAgICAgIDxpbnB1dAogICAgICAgICAgICB0eXBlPSJ0ZXh0IgogICAgICAgICAgICBjbGFzcz0iZm9ybS1jb250cm9sIgogICAgICAgICAgICBpZD0icGhvbmUiCiAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXmiYvmnLrlj7ciCiAgICAgICAgICAvPgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KCiAgICAgIDwhLS0g5a+G56CB6L6T5YWl5qGGIC0tPgogICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgPGxhYmVsIGZvcj0icGFzc3dvcmQiIGNsYXNzPSJjb2wtc20tMiBjb2wtZm9ybS1sYWJlbCI+5a+G56CBPC9sYWJlbD4KICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tMTAiPgogICAgICAgICAgPGlucHV0CiAgICAgICAgICAgIHR5cGU9InBhc3N3b3JkIgogICAgICAgICAgICBjbGFzcz0iZm9ybS1jb250cm9sIgogICAgICAgICAgICBpZD0icGFzc3dvcmQiCiAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXlr4bnoIEiCiAgICAgICAgICAvPgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KCiAgICAgIDwhLS0g5Zyw5Yy656CB6L6T5YWl5qGGIC0tPgogICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgPGxhYmVsIGZvcj0iY29kZSIgY2xhc3M9ImNvbC1zbS0yIGNvbC1mb3JtLWxhYmVsIj7lnLDljLrnoIE8L2xhYmVsPgogICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0xMCI+CiAgICAgICAgICA8aW5wdXQKICAgICAgICAgICAgdHlwZT0idGV4dCIKICAgICAgICAgICAgY2xhc3M9ImZvcm0tY29udHJvbCIKICAgICAgICAgICAgaWQ9ImNvZGUiCiAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXlnLDljLrnoIEiCiAgICAgICAgICAgIHZhbHVlPSI4NiIKICAgICAgICAgIC8+CiAgICAgICAgPC9kaXY+CiAgICAgIDwvZGl2PgoKICAgICAgPCEtLSDlh6Tlh7Dnp4Dku6TniYwgLS0+CiAgICAgIDxkaXYgY2xhc3M9Im1iLTMgcm93Ij4KICAgICAgICA8bGFiZWwgZm9yPSJmZW5nX3Nob3dfdG9rZW4iIGNsYXNzPSJjb2wtc20tMiBjb2wtZm9ybS1sYWJlbCIKICAgICAgICAgID7lh6Tlh7Dnp4Dku6TniYw8L2xhYmVsCiAgICAgICAgPgogICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0xMCI+CiAgICAgICAgICA8dGV4dGFyZWEgY2xhc3M9ImZvcm0tY29udHJvbCIgaWQ9ImZlbmdfc2hvd190b2tlbiIgcm93cz0iNCI+ClVEU1dhbjNnUE5sUlJScnFlMWhqUzRFU2gvUjFSaG9ZMXBFdlFIZW1oN2lsSlhXVHFDR21Iblp0bjdGOXBTcmE2MTc2Q1dxTlo1WUZjZm1hdXhLczk4ajNEdGUwNXhUTm5rdkZQV1BYV0Z3TlhKeXcyc2hkdWVHT241azJwWW82eU0zN3pTRzZ2WGRhYzlXVU1rQVcyUVdJTWM4a0lubnd6ck9VME14VGNnYzwvdGV4dGFyZWEKICAgICAgICAgID4KICAgICAgICA8L2Rpdj4KICAgICAgPC9kaXY+CiAgICA8L2Rpdj4KCiAgICA8IS0tIOaPkOS6pOaMiemSriAtLT4KICAgIDxkaXYgY2xhc3M9InRleHQtY2VudGVyIj4KICAgICAgPGJ1dHRvbiBvbmNsaWNrPSJoYW5kbGVBdXRoKCkiIHR5cGU9ImJ1dHRvbiIgY2xhc3M9ImJ0biBidG4tcHJpbWFyeSI+CiAgICAgICAg55m75b2V6I635Y+WIHRva2VuCiAgICAgIDwvYnV0dG9uPgogICAgPC9kaXY+CgogICAgPCEtLSDnu5Pmnpzlj43ppoggLS0+CiAgICA8ZGl2IGNsYXNzPSJyZXN1bHQtd3JhcCI+CiAgICAgIDxkaXYgY2xhc3M9Im1lc3NhZ2UiPuivt+Whq+WFpeiHquW3seeahOWHpOWHsOengOi0puWPt+S/oeaBr+WQjueCueWHu+eZu+W9leiOt+WPliB0b2tlbjwvZGl2PgogICAgICA8ZGl2IGNsYXNzPSJhdXRoLXRva2VuIj4KICAgICAgICA8dGV4dGFyZWEKICAgICAgICAgIGNsYXNzPSJmb3JtLWNvbnRyb2wiCiAgICAgICAgICBkaXNhYmxlZAogICAgICAgICAgaWQ9ImZlbmdfYXV0aF90b2tlbiIKICAgICAgICAgIHBsYWNlaG9sZGVyPSLmjojmnYPkv6Hmga/lnKjov5nph4zlsZXnpLoiCiAgICAgICAgPjwvdGV4dGFyZWE+CiAgICAgIDwvZGl2PgogICAgPC9kaXY+CiAgPC9ib2R5PgoKICA8c2NyaXB0PgogICAgY29uc3QgaGFuZGxlQXV0aCA9ICgpID0+IHsKICAgICAgY29uc3QgcGhvbmUgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIjcGhvbmUiKS52YWx1ZTsKICAgICAgY29uc3QgcGFzc3dvcmQgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIjcGFzc3dvcmQiKS52YWx1ZTsKICAgICAgY29uc3QgY29kZSA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoIiNjb2RlIikudmFsdWU7CiAgICAgIGNvbnN0IGZlbmdTaG93VG9rZW4gPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIjZmVuZ19zaG93X3Rva2VuIikudmFsdWU7CiAgICAgIGNvbnN0IGRhdGEgPSB7cGhvbmUsIHBhc3N3b3JkLCBjb2RlLCB0aWNrZXQ6IGZlbmdTaG93VG9rZW59OwoKICAgICAgY29uc3QgbWVzc2FnZUVsbSA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoIi5yZXN1bHQtd3JhcCAubWVzc2FnZSIpOwogICAgICBtZXNzYWdlRWxtLmlubmVySFRNTCA9ICLmraPlnKjmjojmnYPkuK0sIOivt+eojeWQji4uLiI7CiAgICAgIGNvbnN0IHRva2VuRWxtID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvcigiLnJlc3VsdC13cmFwICNmZW5nX2F1dGhfdG9rZW4iKTsKCiAgICAgIC8vIOWPkei1tyBQT1NUIOivt+axggogICAgICBmZXRjaCgiaHR0cDovL20uZmVuZ3Nob3dzLmNvbS91c2VyL29hdXRoL2xvZ2luIiwgewogICAgICAgIG1ldGhvZDogIlBPU1QiLCAvLyDkvb/nlKggUE9TVCDmlrnms5UKICAgICAgICBoZWFkZXJzOiB7CiAgICAgICAgICAiQ29udGVudC1UeXBlIjogImFwcGxpY2F0aW9uL2pzb24iLCAvLyDmjIflrpror7fmsYLkvZPnsbvlnovkuLogSlNPTgogICAgICAgIH0sCiAgICAgICAgYm9keTogSlNPTi5zdHJpbmdpZnkoZGF0YSksIC8vIOWwhuWvueixoei9rOaNouS4uiBKU09OIOWtl+espuS4sgogICAgICB9KQogICAgICAgIC50aGVuKChyZXNwb25zZSkgPT4gewogICAgICAgICAgaWYgKCFyZXNwb25zZS5vaykgewogICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYEhUVFAgZXJyb3IhIFN0YXR1czogJHtyZXNwb25zZS5zdGF0dXN9YCk7CiAgICAgICAgICB9CiAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpOyAvLyDop6PmnpAgSlNPTiDlk43lupTkvZMKICAgICAgICB9KQogICAgICAgIC50aGVuKChkYXRhKSA9PiB7CiAgICAgICAgICBpZiAoZGF0YS5zdGF0dXMgIT09ICcwJykgewogICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYOaOiOadg+Wksei0pTogJHtkYXRhLm1lc3NhZ2V9YCk7CiAgICAgICAgICB9CiAgICAgICAgICBjb25zdCB0b2tlbiA9IGRhdGE/LmRhdGE/LnRva2VuOwogICAgICAgICAgbWVzc2FnZUVsbS5pbm5lckhUTUwgPSAn5o6I5p2D5oiQ5Yqf77yM6K+35bCG6I635Y+W5Yiw55qE5YC86K6+572u5Yiw56iL5bqP55qE546v5aKD5Y+Y6YeP5Lit77yMa2V5OiBmZW5nX3Rva2VuJzsKICAgICAgICAgIHRva2VuRWxtLnZhbHVlID0gdG9rZW47CiAgICAgICAgfSkKICAgICAgICAuY2F0Y2goKGVycm9yKSA9PiB7CiAgICAgICAgICBtZXNzYWdlRWxtLmlubmVySFRNTCA9IGVycm9yOwogICAgICAgICAgc2V0VGltZW91dCgoKSA9PiBtZXNzYWdlRWxtLmlubmVySFRNTCA9ICfor7floavlhaXoh6rlt7HnmoTlh6Tlh7Dnp4DotKblj7fkv6Hmga/lkI7ngrnlh7vnmbvlvZXojrflj5YgdG9rZW4nLCA1MDAwKTsKICAgICAgICB9KTsKICAgIH07CiAgPC9zY3JpcHQ+CjwvaHRtbD4K`

	// ConfigPageHtml 配置页地址
	ConfigPageHtml = `PCFET0NUWVBFIGh0bWw+CjxodG1sPgogIDxoZWFkPgogICAgPG1ldGEgY2hhcnNldD0iVVRGLTgiIC8+CiAgICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBpbml0aWFsLXNjYWxlPTEuMCIgLz4KICAgIDx0aXRsZT5saXZlLXNlcnZlciDphY3nva7pobU8L3RpdGxlPgogICAgPCEtLSDlvJXlhaUgQm9vdHN0cmFwIENTUyAtLT4KICAgIDxsaW5rCiAgICAgIGhyZWY9Imh0dHBzOi8vZmFzdGx5LmpzZGVsaXZyLm5ldC9ucG0vYm9vdHN0cmFwQDUuMy4wL2Rpc3QvY3NzL2Jvb3RzdHJhcC5taW4uY3NzIgogICAgICByZWw9InN0eWxlc2hlZXQiCiAgICAvPgogICAgPHN0eWxlPgogICAgICAucGFnZS1jb250YWluZXIgewogICAgICAgIHBhZGRpbmc6IDIwcHggMTBweDsKICAgICAgfQogICAgICAucGFnZS1jb250YWluZXIgaW5wdXQsCiAgICAgIC5wYWdlLWNvbnRhaW5lciB0ZXh0YXJlYSwKICAgICAgLnBhZ2UtY29udGFpbmVyIC5idXR0b24tZ3JvdXAgewogICAgICAgIG1hcmdpbjogOHB4IDA7CiAgICAgIH0KICAgICAgLmVudi13cmFwIC5lbnYtdmFsdWUgewogICAgICAgIG92ZXJmbG93LXdyYXA6IGJyZWFrLXdvcmQ7CiAgICAgICAgd29yZC1icmVhazogYnJlYWstYWxsOwogICAgICB9CiAgICA8L3N0eWxlPgogIDwvaGVhZD4KICA8Ym9keT4KICAgIDxkaXYgY2xhc3M9InBhZ2UtY29udGFpbmVyIj4KICAgICAgPGgyPmxpdmUtc2VydmVyIOmFjee9rumhtTwvaDI+CgogICAgICA8IS0tIOWxleekuueUqOaIt+S8oOmAkueahOeoi+W6j+WvhumSpSAtLT4KICAgICAgPGRpdiBjbGFzcz0ibWItMyByb3cgc2VjcmV0LXdyYXAiPgogICAgICAgIDxsYWJlbCBjbGFzcz0iY29sLXNtLTIgY29sLWZvcm0tbGFiZWwiPueoi+W6j+WvhumSpe+8mjwvbGFiZWw+CiAgICAgICAgPGRpdiBjbGFzcz0iY29sLXNtLTEwIj4KICAgICAgICAgIDxpbnB1dCBjbGFzcz0iZm9ybS1jb250cm9sIHNlY3JldCIgdHlwZT0idGV4dCIgZGlzYWJsZWQgLz4KICAgICAgICA8L2Rpdj4KICAgICAgPC9kaXY+CgogICAgICA8IS0tIOeOr+Wig+WPmOmHjyAtLT4KICAgICAgPGRpdiBjbGFzcz0iZW52LXdyYXAiPgogICAgICAgIDxoMz7njq/looPlj5jph488L2gzPgogICAgICAgIDxkaXYgY2xhc3M9Im1iLTMgcm93Ij4KICAgICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS02Ij4KICAgICAgICAgICAgPGlucHV0CiAgICAgICAgICAgICAgY2xhc3M9ImZvcm0tY29udHJvbCBlbnYta2V5IgogICAgICAgICAgICAgIHR5cGU9InRleHQiCiAgICAgICAgICAgICAgcGxhY2Vob2xkZXI9IuWcqOatpOi+k+WFpSBrZXksIOWmgjogZmVuZ190b2tlbiIKICAgICAgICAgICAgLz4KICAgICAgICAgIDwvZGl2PgogICAgICAgICAgPGRpdiBjbGFzcz0iY29sLXNtLTYgYnV0dG9uLWdyb3VwIj4KICAgICAgICAgICAgPGJ1dHRvbiB0eXBlPSJidXR0b24iIGNsYXNzPSJidG4gYnRuLXN1Y2Nlc3MiIG9uY2xpY2s9ImdldEVudigpIj4KICAgICAgICAgICAgICDmn6Ug6K+iCiAgICAgICAgICAgIDwvYnV0dG9uPgogICAgICAgICAgICA8YnV0dG9uIHR5cGU9ImJ1dHRvbiIgY2xhc3M9ImJ0biBidG4tcHJpbWFyeSIgb25jbGljaz0ic2V0RW52KCkiPgogICAgICAgICAgICAgIOa3uyDliqAKICAgICAgICAgICAgPC9idXR0b24+CiAgICAgICAgICAgIDxidXR0b24gdHlwZT0iYnV0dG9uIiBjbGFzcz0iYnRuIGJ0bi1kYW5nZXIiIG9uY2xpY2s9ImRlbEVudigpIj4KICAgICAgICAgICAgICDliKAg6ZmkCiAgICAgICAgICAgIDwvYnV0dG9uPgogICAgICAgICAgPC9kaXY+CiAgICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tMTIiPgogICAgICAgICAgICA8dGV4dGFyZWEKICAgICAgICAgICAgICBjbGFzcz0iZm9ybS1jb250cm9sIGVudi12YWx1ZSIKICAgICAgICAgICAgICByb3dzPSI2IgogICAgICAgICAgICAgIHN0eWxlPSJvdmVyZmxvdzogaGlkZGVuOyByZXNpemU6IG5vbmUiCiAgICAgICAgICAgICAgcGxhY2Vob2xkZXI9IuWcqOatpOi+k+WFpeaDs+abtOaWsOeahCB2YWx1ZSwg6YCa6L+H5p+l6K+i5b6X5Yiw55qE5YC85Lmf5Lya5pi+56S65Zyo5q2kIgogICAgICAgICAgICA+PC90ZXh0YXJlYT4KICAgICAgICAgIDwvZGl2PgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KCiAgICAgIDwhLS0g6buR5ZCN5Y2VIC0tPgogICAgICA8ZGl2IGNsYXNzPSJibGFjay1pcC13cmFwIj4KICAgICAgICA8aDM+6buR5ZCN5Y2VPC9oMz4KICAgICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tMTAiPgogICAgICAgICAgICA8aW5wdXQKICAgICAgICAgICAgICBjbGFzcz0iZm9ybS1jb250cm9sIGJsYWNrLWlwIgogICAgICAgICAgICAgIHR5cGU9InRleHQiCiAgICAgICAgICAgICAgcGxhY2Vob2xkZXI9Iuivt+i+k+WFpeimgeWKoOWFpem7keWQjeWNleeahCBpcCwg5aaCOiAxLjEuMS4xIgogICAgICAgICAgICAvPgogICAgICAgICAgPC9kaXY+CiAgICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tMiBidXR0b24tZ3JvdXAiPgogICAgICAgICAgICA8YnV0dG9uCiAgICAgICAgICAgICAgdHlwZT0iYnV0dG9uIgogICAgICAgICAgICAgIGNsYXNzPSJidG4gYnRuLXByaW1hcnkiCiAgICAgICAgICAgICAgb25jbGljaz0ic2V0QmxhY2tJcCgpIgogICAgICAgICAgICA+CiAgICAgICAgICAgICAg5re7IOWKoAogICAgICAgICAgICA8L2J1dHRvbj4KICAgICAgICAgIDwvZGl2PgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KCiAgICAgIDwhLS0g5Zyw5Z+f55m95ZCN5Y2VIC0tPgogICAgICA8ZGl2IGNsYXNzPSJ3aGl0ZS1hcmVhLXdyYXAiPgogICAgICAgIDxoMz7lnLDln5/nmb3lkI3ljZU8L2gzPgogICAgICAgIDxkaXYgY2xhc3M9Im1iLTMgcm93Ij4KICAgICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS04Ij4KICAgICAgICAgICAgPGlucHV0CiAgICAgICAgICAgICAgY2xhc3M9ImZvcm0tY29udHJvbCB3aGl0ZS1hcmVhIgogICAgICAgICAgICAgIHR5cGU9InRleHQiCiAgICAgICAgICAgICAgcGxhY2Vob2xkZXI9Iuivt+i+k+WFpeimgeWKoOWFpeeZveWQjeWNleeahOWcsOWfnywg5aaCOiDlub/kuJwv5L2b5bGxL+WNl+a1tyIKICAgICAgICAgICAgLz4KICAgICAgICAgIDwvZGl2PgogICAgICAgICAgPGRpdiBjbGFzcz0iY29sLXNtLTQgYnV0dG9uLWdyb3VwIj4KICAgICAgICAgICAgPGJ1dHRvbgogICAgICAgICAgICAgIHR5cGU9ImJ1dHRvbiIKICAgICAgICAgICAgICBjbGFzcz0iYnRuIGJ0bi1wcmltYXJ5IgogICAgICAgICAgICAgIG9uY2xpY2s9InNldFdoaXRlQXJlYSgnc2V0JykiCiAgICAgICAgICAgID4KICAgICAgICAgICAgICDmt7sg5YqgCiAgICAgICAgICAgIDwvYnV0dG9uPgogICAgICAgICAgICA8YnV0dG9uCiAgICAgICAgICAgICAgdHlwZT0iYnV0dG9uIgogICAgICAgICAgICAgIGNsYXNzPSJidG4gYnRuLWRhbmdlciIKICAgICAgICAgICAgICBvbmNsaWNrPSJzZXRXaGl0ZUFyZWEoJ2RlbCcpIgogICAgICAgICAgICA+CiAgICAgICAgICAgICAg56e7IOmZpAogICAgICAgICAgICA8L2J1dHRvbj4KICAgICAgICAgIDwvZGl2PgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KICAgIDwvZGl2PgogIDwvYm9keT4KCiAgPHNjcmlwdD4KICAgIC8vIOiOt+WPluWfuuacrOWPguaVsAogICAgY29uc3QgZ2V0QmFzZVBhcmFtcyA9ICgpID0+IHsKICAgICAgY29uc3Qgc2VjcmV0ID0gbmV3IFVSTCh3aW5kb3cubG9jYXRpb24uaHJlZikuc2VhcmNoUGFyYW1zLmdldCgic2VjcmV0Iik7CiAgICAgIHJldHVybiB7CiAgICAgICAgc2VjcmV0LCAvLyDnqIvluo/lr4bpkqUKICAgICAgICBob3N0OiAiIiwgLy8g5Y+R6YCB6K+35rGC55qE5Li75py65Zyw5Z2ACiAgICAgIH07CiAgICB9OwoKICAgIC8vIOagueaNruaWh+acrOahhuWGheWuueiHquWKqOiwg+aVtOeOr+Wig+WPmOmHj+ahhueahOmrmOW6pgogICAgY29uc3QgYXV0b1Jlc2l6ZUVudlZhbHVlVGV4dGFyZWEgPSAoaW5pdEV2ZW50ID0gZmFsc2UpID0+IHsKICAgICAgY29uc3QgdGV4dGFyZWEgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIuZW52LXdyYXAgLmVudi12YWx1ZSIpOwoKICAgICAgLy8g5qC45b+D6LCD5pW05Ye95pWwCiAgICAgIGNvbnN0IGF1dG9SZXNpemUgPSAoKSA9PiB7CiAgICAgICAgdGV4dGFyZWEuc3R5bGUuaGVpZ2h0ID0gImF1dG8iOwogICAgICAgIHRleHRhcmVhLnN0eWxlLmhlaWdodCA9IGAke3RleHRhcmVhLnNjcm9sbEhlaWdodH1weGA7CiAgICAgIH07CgogICAgICBpZiAoaW5pdEV2ZW50KSB7CiAgICAgICAgLy8g6L6T5YWl5LqL5Lu255uR5ZCsCiAgICAgICAgdGV4dGFyZWEuYWRkRXZlbnRMaXN0ZW5lcigiaW5wdXQiLCBhdXRvUmVzaXplKTsKCiAgICAgICAgLy8g56qX5Y+j5bC65a+45Y+Y5YyW55uR5ZCs77yI5bim6Ziy5oqW77yJCiAgICAgICAgbGV0IHJlc2l6ZVRpbWVyOwogICAgICAgIHdpbmRvdy5hZGRFdmVudExpc3RlbmVyKCJyZXNpemUiLCAoKSA9PiB7CiAgICAgICAgICBjbGVhclRpbWVvdXQocmVzaXplVGltZXIpOwogICAgICAgICAgcmVzaXplVGltZXIgPSBzZXRUaW1lb3V0KGF1dG9SZXNpemUsIDE1MCk7CiAgICAgICAgfSk7CiAgICAgIH0KCiAgICAgIC8vIOWIneWni+WMluaJp+ihjAogICAgICBhdXRvUmVzaXplKCk7CiAgICB9OwoKICAgIHdpbmRvdy5vbmxvYWQgPSAoKSA9PiB7CiAgICAgIC8vIOiuvue9rueoi+W6j+WvhumSpQogICAgICBjb25zdCB7IHNlY3JldCB9ID0gZ2V0QmFzZVBhcmFtcygpOwogICAgICBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIucGFnZS1jb250YWluZXIgLnNlY3JldCIpLnZhbHVlID0gc2VjcmV0OwoKICAgICAgYXV0b1Jlc2l6ZUVudlZhbHVlVGV4dGFyZWEodHJ1ZSk7CiAgICB9OwoKICAgIC8vIOiOt+WPlui+k+WFpeahhueahOWAvAogICAgY29uc3QgZ2V0SW5wdXRWYWx1ZSA9IChzZWxlY3RvciA9ICIiKSA9PiB7CiAgICAgIHJldHVybiAoZG9jdW1lbnQucXVlcnlTZWxlY3RvcihzZWxlY3RvcikudmFsdWUgfHwgIiIpLnRyaW0oKTsKICAgIH07CgogICAgLy8g6K6+572u6L6T5YWl5qGG55qE5YC8CiAgICBjb25zdCBzZXRJbnB1dFZhbHVlID0gKHNlbGVjdG9yID0gIiIsIHZhbHVlID0gIiIpID0+IHsKICAgICAgY29uc3QgdGFyZ2V0ID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvcihzZWxlY3Rvcik7CiAgICAgIGlmICghdGFyZ2V0KSB7CiAgICAgICAgcmV0dXJuOwogICAgICB9CiAgICAgIHRhcmdldC52YWx1ZSA9IHZhbHVlOwogICAgfTsKCiAgICAvLyBHRVQg6K+35rGC5oyH5a6a5Zyw5Z2ACiAgICBjb25zdCBmZXRjaFVybCA9ICgKICAgICAgdXJsID0gIiIsCiAgICAgIHN1Y2Nlc3NGdW5jID0gKCkgPT4ge30sCiAgICAgIGVycm9yRnVuYyA9IChlcnIpID0+IHt9CiAgICApID0+IHsKICAgICAgZmV0Y2godXJsKQogICAgICAgIC50aGVuKChyZXMpID0+IHsKICAgICAgICAgIGlmIChyZXMub2spIHsKICAgICAgICAgICAgcmV0dXJuIHN1Y2Nlc3NGdW5jKCk7CiAgICAgICAgICB9CiAgICAgICAgICBlcnJvckZ1bmMocmVzLnN0YXR1c1RleHQpOwogICAgICAgIH0pCiAgICAgICAgLmNhdGNoKChlcnIpID0+IHsKICAgICAgICAgIGVycm9yRnVuYyhlcnIpOwogICAgICAgIH0pOwogICAgfTsKCiAgICAvLyDmn6Xor6Lnjq/looPlj5jph48KICAgIGNvbnN0IGdldEVudiA9ICgpID0+IHsKICAgICAgY29uc3QgZW52S2V5ID0gZ2V0SW5wdXRWYWx1ZSgiLmVudi13cmFwIC5lbnYta2V5Iik7CiAgICAgIGlmIChlbnZLZXkgPT09ICIiKSB7CiAgICAgICAgcmV0dXJuIGFsZXJ0KCLovpPlhaXlj4LmlbDkuI3og73kuLrnqboiKTsKICAgICAgfQoKICAgICAgY29uc3QgeyBzZWNyZXQsIGhvc3QgfSA9IGdldEJhc2VQYXJhbXMoKTsKICAgICAgY29uc3QgdXJsID0gYCR7aG9zdH0vZW52P2tleT0ke2VudktleX0mc2VjcmV0PSR7c2VjcmV0fWA7CiAgICAgIGZldGNoKHVybCkKICAgICAgICAudGhlbigocmVzKSA9PiB7CiAgICAgICAgICBpZiAoIXJlcy5vaykgewogICAgICAgICAgICByZXR1cm4gYWxlcnQoYOafpeivouWksei0pTogJHtyZXMuc3RhdHVzVGV4dH1gKTsKICAgICAgICAgIH0KICAgICAgICAgIHJldHVybiByZXMudGV4dCgpOwogICAgICAgIH0pCiAgICAgICAgLnRoZW4oKHJlc3BUZXh0KSA9PiB7CiAgICAgICAgICBzZXRJbnB1dFZhbHVlKCIuZW52LXdyYXAgLmVudi12YWx1ZSIsIHJlc3BUZXh0KTsKICAgICAgICAgIGF1dG9SZXNpemVFbnZWYWx1ZVRleHRhcmVhKGZhbHNlKTsKICAgICAgICB9KQogICAgICAgIC5jYXRjaCgoZXJyKSA9PiB7CiAgICAgICAgICBhbGVydChg5p+l6K+i5aSx6LSlOiAke2Vycn1gKTsKICAgICAgICB9KTsKICAgIH07CgogICAgLy8g6K6+572u546v5aKD5Y+Y6YePCiAgICBjb25zdCBzZXRFbnYgPSAoKSA9PiB7CiAgICAgIGNvbnN0IGVudktleSA9IGdldElucHV0VmFsdWUoIi5lbnYtd3JhcCAuZW52LWtleSIpOwogICAgICBjb25zdCBlbnZWYWx1ZSA9IGdldElucHV0VmFsdWUoIi5lbnYtd3JhcCAuZW52LXZhbHVlIik7CiAgICAgIGlmIChlbnZLZXkgPT09ICIiIHx8IGVudlZhbHVlID09PSAiIikgewogICAgICAgIHJldHVybiBhbGVydCgi6L6T5YWl5Y+C5pWw5LiN6IO95Li656m6Iik7CiAgICAgIH0KCiAgICAgIGNvbnN0IHsgc2VjcmV0LCBob3N0IH0gPSBnZXRCYXNlUGFyYW1zKCk7CiAgICAgIGNvbnN0IHVybCA9IGAke2hvc3R9L2Vudj9zZWNyZXQ9JHtzZWNyZXR9YDsKICAgICAgY29uc3QgZm9ybURhdGEgPSBuZXcgRm9ybURhdGEoKTsKICAgICAgZm9ybURhdGEuYXBwZW5kKCJrZXkiLCBlbnZLZXkpOwogICAgICBmb3JtRGF0YS5hcHBlbmQoInZhbHVlIiwgZW52VmFsdWUpOwoKICAgICAgZmV0Y2godXJsLCB7CiAgICAgICAgbWV0aG9kOiAiUE9TVCIsCiAgICAgICAgYm9keTogZm9ybURhdGEsCiAgICAgIH0pCiAgICAgICAgLnRoZW4oKHJlcykgPT4gewogICAgICAgICAgaWYgKHJlcy5vaykgewogICAgICAgICAgICByZXR1cm4gYWxlcnQoIua3u+WKoOaIkOWKnyIpOwogICAgICAgICAgfQogICAgICAgICAgYWxlcnQoYOa3u+WKoOWksei0pTogJHtyZXMuc3RhdHVzVGV4dH1gKTsKICAgICAgICB9KQogICAgICAgIC5jYXRjaCgoZXJyKSA9PiB7CiAgICAgICAgICBhbGVydChg5re75Yqg5aSx6LSlOiAke2Vycn1gKTsKICAgICAgICB9KTsKICAgIH07CgogICAgLy8g5Yig6Zmk546v5aKD5Y+Y6YePCiAgICBjb25zdCBkZWxFbnYgPSAoKSA9PiB7CiAgICAgIGNvbnN0IGVudktleSA9IGdldElucHV0VmFsdWUoIi5lbnYtd3JhcCAuZW52LWtleSIpOwogICAgICBpZiAoZW52S2V5ID09PSAiIikgewogICAgICAgIHJldHVybiBhbGVydCgi6L6T5YWl5Y+C5pWw5LiN6IO95Li656m6Iik7CiAgICAgIH0KCiAgICAgIGNvbnN0IHsgc2VjcmV0LCBob3N0IH0gPSBnZXRCYXNlUGFyYW1zKCk7CiAgICAgIGNvbnN0IHVybCA9IGAke2hvc3R9L2Vudj9zZWNyZXQ9JHtzZWNyZXR9YDsKICAgICAgY29uc3QgZm9ybURhdGEgPSBuZXcgRm9ybURhdGEoKTsKICAgICAgZm9ybURhdGEuYXBwZW5kKCJrZXkiLCBlbnZLZXkpOwoKICAgICAgZmV0Y2godXJsLCB7CiAgICAgICAgbWV0aG9kOiAiREVMRVRFIiwKICAgICAgICBib2R5OiBmb3JtRGF0YSwKICAgICAgfSkKICAgICAgICAudGhlbigocmVzKSA9PiB7CiAgICAgICAgICBpZiAocmVzLm9rKSB7CiAgICAgICAgICAgIHNldElucHV0VmFsdWUoIi5lbnYtd3JhcCAuZW52LXZhbHVlIiwgIiIpOwogICAgICAgICAgICByZXR1cm4gYWxlcnQoIuWIoOmZpOaIkOWKnyIpOwogICAgICAgICAgfQogICAgICAgICAgYWxlcnQoYOWIoOmZpOWksei0pTogJHtyZXMuc3RhdHVzVGV4dH1gKTsKICAgICAgICB9KQogICAgICAgIC5jYXRjaCgoZXJyKSA9PiB7CiAgICAgICAgICBhbGVydChg5Yig6Zmk5aSx6LSlOiAke2Vycn1gKTsKICAgICAgICB9KTsKICAgIH07CgogICAgLy8g6K6+572u6buR5ZCN5Y2VIGlwCiAgICBjb25zdCBzZXRCbGFja0lwID0gKCkgPT4gewogICAgICBjb25zdCBibGFja0lwID0gZ2V0SW5wdXRWYWx1ZSgiLmJsYWNrLWlwLXdyYXAgLmJsYWNrLWlwIik7CiAgICAgIGlmIChibGFja0lwID09PSAiIikgewogICAgICAgIHJldHVybiBhbGVydCgi5Y+C5pWw5LiN6IO95Li656m6Iik7CiAgICAgIH0KICAgICAgY29uc3QgeyBzZWNyZXQsIGhvc3QgfSA9IGdldEJhc2VQYXJhbXMoKTsKICAgICAgY29uc3QgdXJsID0gYCR7aG9zdH0vYmxhY2tfaXA/aXA9JHtibGFja0lwfSZzZWNyZXQ9JHtzZWNyZXR9YDsKICAgICAgZmV0Y2hVcmwoCiAgICAgICAgdXJsLAogICAgICAgICgpID0+IGFsZXJ0KCLmt7vliqDmiJDlip8iKSwKICAgICAgICAoZXJyKSA9PiBhbGVydChg5re75Yqg5aSx6LSlOiAke2Vycn1gKQogICAgICApOwogICAgfTsKCiAgICAvLyDorr7nva7lnLDln5/nmb3lkI3ljZUKICAgIGNvbnN0IHNldFdoaXRlQXJlYSA9IChhY3Rpb24gPSAic2V0IikgPT4gewogICAgICBjb25zdCB3aGl0ZUFyZWEgPSBnZXRJbnB1dFZhbHVlKCIud2hpdGUtYXJlYS13cmFwIC53aGl0ZS1hcmVhIik7CiAgICAgIGlmICh3aGl0ZUFyZWEgPT09ICIiKSB7CiAgICAgICAgcmV0dXJuIGFsZXJ0KCLovpPlhaXlj4LmlbDkuI3og73kuLrnqboiKTsKICAgICAgfQogICAgICBjb25zdCBhY3Rpb25UaXAgPSBhY3Rpb24gPT09ICJzZXQiID8gIua3u+WKoCIgOiAi56e76ZmkIjsKICAgICAgY29uc3QgeyBzZWNyZXQsIGhvc3QgfSA9IGdldEJhc2VQYXJhbXMoKTsKICAgICAgY29uc3QgdXJsID0gYCR7aG9zdH0vd2hpdGVfYXJlYS8ke2FjdGlvbn0/YXJlYT0ke3doaXRlQXJlYX0mc2VjcmV0PSR7c2VjcmV0fWA7CiAgICAgIGZldGNoVXJsKAogICAgICAgIHVybCwKICAgICAgICAoKSA9PiBhbGVydChgJHthY3Rpb25UaXB95oiQ5YqfYCksCiAgICAgICAgKGVycikgPT4gYWxlcnQoYCR7YWN0aW9uVGlwfeWksei0pTogJHtlcnJ9YCkKICAgICAgKTsKICAgIH07CiAgPC9zY3JpcHQ+CjwvaHRtbD4K`
)

const (
	Gin_IpAddrInfoKey = "ip_addr_info" // 在 gin 上下文中存放 IP 地址信息的键
)
