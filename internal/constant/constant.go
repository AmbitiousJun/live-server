package constant

const (
	Version = "v1.16.4" // 版本号
)

const (
	Dir_DataRoot = "data" // 数据存放根目录
	// Dir_DataRoot = "/Users/ambitious/Desktop/code/go/live-server/data" // 数据存放根目录, 开发环境
)

const (

	// HelpDocHtmlTemplate 帮助文档的 html 模板
	HelpDocHtmlTemplate = `PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KICA8aGVhZD4KICAgIDxtZXRhIGNoYXJzZXQ9IlVURi04IiAvPgogICAgPG1ldGEgbmFtZT0idmlld3BvcnQiIGNvbnRlbnQ9IndpZHRoPWRldmljZS13aWR0aCwgaW5pdGlhbC1zY2FsZT0xLjAiIC8+CiAgICA8dGl0bGU+bGl2ZS1zZXJ2ZXIg5biu5Yqp5paH5qGjPC90aXRsZT4KICA8L2hlYWQ+CiAgPGJvZHk+CiAgICA8ZGl2IGNsYXNzPSJjb250YWluZXIiPiR7ZG9jQ29udGVudH08L2Rpdj4KCiAgICA8c2NyaXB0PgogICAgICB3aW5kb3cub25sb2FkID0gKCkgPT4gewogICAgICAgIGNvbnN0IGNvbnRhaW5lciA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoJy5jb250YWluZXInKTsKICAgICAgICBjb250YWluZXIuaW5uZXJIVE1MID0gY29udGFpbmVyLmlubmVySFRNTC5yZXBsYWNlKC9cJFx7Y2xpZW50T3JpZ2luXH0vZywgbG9jYXRpb24ub3JpZ2luKTsKICAgICAgfQogICAgPC9zY3JpcHQ+CiAgPC9ib2R5Pgo8L2h0bWw+Cg==`

	// FengAuthHtml 凤凰秀授权地址
	FengAuthHtml = `PCFET0NUWVBFIGh0bWw+CjxodG1sPgogIDxoZWFkPgogICAgPG1ldGEgY2hhcnNldD0iVVRGLTgiIC8+CiAgICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBpbml0aWFsLXNjYWxlPTEuMCIgLz4KICAgIDx0aXRsZT7lh6Tlh7Dnp4DnmbvlvZXmjojmnYM8L3RpdGxlPgogICAgPCEtLSDlvJXlhaUgQm9vdHN0cmFwIENTUyAtLT4KICAgIDxsaW5rCiAgICAgIGhyZWY9Imh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9ucG0vYm9vdHN0cmFwQDUuMy4wL2Rpc3QvY3NzL2Jvb3RzdHJhcC5taW4uY3NzIgogICAgICByZWw9InN0eWxlc2hlZXQiCiAgICAvPgogICAgPHN0eWxlPgogICAgICAucmVzdWx0LXdyYXAgewogICAgICAgIG1hcmdpbi10b3A6IDQwcHg7CiAgICAgICAgcGFkZGluZzogMjBweDsKICAgICAgfQogICAgICAucmVzdWx0LXdyYXAgLm1lc3NhZ2UgewogICAgICAgIHBhZGRpbmc6IDEwcHg7CiAgICAgIH0KICAgIDwvc3R5bGU+CiAgPC9oZWFkPgogIDxib2R5PgogICAgPGRpdiBjbGFzcz0iY29udGFpbmVyIG10LTQiPgogICAgICA8IS0tIOaJi+acuuWPt+i+k+WFpeahhiAtLT4KICAgICAgPGRpdiBjbGFzcz0ibWItMyByb3ciPgogICAgICAgIDxsYWJlbCBmb3I9InBob25lIiBjbGFzcz0iY29sLXNtLTIgY29sLWZvcm0tbGFiZWwiPuaJi+acuuWPtzwvbGFiZWw+CiAgICAgICAgPGRpdiBjbGFzcz0iY29sLXNtLTEwIj4KICAgICAgICAgIDxpbnB1dAogICAgICAgICAgICB0eXBlPSJ0ZXh0IgogICAgICAgICAgICBjbGFzcz0iZm9ybS1jb250cm9sIgogICAgICAgICAgICBpZD0icGhvbmUiCiAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXmiYvmnLrlj7ciCiAgICAgICAgICAvPgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KCiAgICAgIDwhLS0g5a+G56CB6L6T5YWl5qGGIC0tPgogICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgPGxhYmVsIGZvcj0icGFzc3dvcmQiIGNsYXNzPSJjb2wtc20tMiBjb2wtZm9ybS1sYWJlbCI+5a+G56CBPC9sYWJlbD4KICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tMTAiPgogICAgICAgICAgPGlucHV0CiAgICAgICAgICAgIHR5cGU9InBhc3N3b3JkIgogICAgICAgICAgICBjbGFzcz0iZm9ybS1jb250cm9sIgogICAgICAgICAgICBpZD0icGFzc3dvcmQiCiAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXlr4bnoIEiCiAgICAgICAgICAvPgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KCiAgICAgIDwhLS0g5Zyw5Yy656CB6L6T5YWl5qGGIC0tPgogICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgPGxhYmVsIGZvcj0iY29kZSIgY2xhc3M9ImNvbC1zbS0yIGNvbC1mb3JtLWxhYmVsIj7lnLDljLrnoIE8L2xhYmVsPgogICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0xMCI+CiAgICAgICAgICA8aW5wdXQKICAgICAgICAgICAgdHlwZT0idGV4dCIKICAgICAgICAgICAgY2xhc3M9ImZvcm0tY29udHJvbCIKICAgICAgICAgICAgaWQ9ImNvZGUiCiAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXlnLDljLrnoIEiCiAgICAgICAgICAgIHZhbHVlPSI4NiIKICAgICAgICAgIC8+CiAgICAgICAgPC9kaXY+CiAgICAgIDwvZGl2PgoKICAgICAgPCEtLSDlh6Tlh7Dnp4Dku6TniYwgLS0+CiAgICAgIDxkaXYgY2xhc3M9Im1iLTMgcm93Ij4KICAgICAgICA8bGFiZWwgZm9yPSJmZW5nX3Nob3dfdG9rZW4iIGNsYXNzPSJjb2wtc20tMiBjb2wtZm9ybS1sYWJlbCIKICAgICAgICAgID7lh6Tlh7Dnp4Dku6TniYw8L2xhYmVsCiAgICAgICAgPgogICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0xMCI+CiAgICAgICAgICA8dGV4dGFyZWEgY2xhc3M9ImZvcm0tY29udHJvbCIgaWQ9ImZlbmdfc2hvd190b2tlbiIgcm93cz0iNCI+ClVEU1dhbjNnUE5sUlJScnFlMWhqUzRFU2gvUjFSaG9ZMXBFdlFIZW1oN2lsSlhXVHFDR21Iblp0bjdGOXBTcmE2MTc2Q1dxTlo1WUZjZm1hdXhLczk4ajNEdGUwNXhUTm5rdkZQV1BYV0Z3TlhKeXcyc2hkdWVHT241azJwWW82eU0zN3pTRzZ2WGRhYzlXVU1rQVcyUVdJTWM4a0lubnd6ck9VME14VGNnYzwvdGV4dGFyZWEKICAgICAgICAgID4KICAgICAgICA8L2Rpdj4KICAgICAgPC9kaXY+CiAgICA8L2Rpdj4KCiAgICA8IS0tIOaPkOS6pOaMiemSriAtLT4KICAgIDxkaXYgY2xhc3M9InRleHQtY2VudGVyIj4KICAgICAgPGJ1dHRvbiBvbmNsaWNrPSJoYW5kbGVBdXRoKCkiIHR5cGU9ImJ1dHRvbiIgY2xhc3M9ImJ0biBidG4tcHJpbWFyeSI+CiAgICAgICAg55m75b2V6I635Y+WIHRva2VuCiAgICAgIDwvYnV0dG9uPgogICAgPC9kaXY+CgogICAgPCEtLSDnu5Pmnpzlj43ppoggLS0+CiAgICA8ZGl2IGNsYXNzPSJyZXN1bHQtd3JhcCI+CiAgICAgIDxkaXYgY2xhc3M9Im1lc3NhZ2UiPuivt+Whq+WFpeiHquW3seeahOWHpOWHsOengOi0puWPt+S/oeaBr+WQjueCueWHu+eZu+W9leiOt+WPliB0b2tlbjwvZGl2PgogICAgICA8ZGl2IGNsYXNzPSJhdXRoLXRva2VuIj4KICAgICAgICA8dGV4dGFyZWEKICAgICAgICAgIGNsYXNzPSJmb3JtLWNvbnRyb2wiCiAgICAgICAgICBkaXNhYmxlZAogICAgICAgICAgaWQ9ImZlbmdfYXV0aF90b2tlbiIKICAgICAgICAgIHBsYWNlaG9sZGVyPSLmjojmnYPkv6Hmga/lnKjov5nph4zlsZXnpLoiCiAgICAgICAgPjwvdGV4dGFyZWE+CiAgICAgIDwvZGl2PgogICAgPC9kaXY+CiAgPC9ib2R5PgoKICA8c2NyaXB0PgogICAgY29uc3QgaGFuZGxlQXV0aCA9ICgpID0+IHsKICAgICAgY29uc3QgcGhvbmUgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIjcGhvbmUiKS52YWx1ZTsKICAgICAgY29uc3QgcGFzc3dvcmQgPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIjcGFzc3dvcmQiKS52YWx1ZTsKICAgICAgY29uc3QgY29kZSA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoIiNjb2RlIikudmFsdWU7CiAgICAgIGNvbnN0IGZlbmdTaG93VG9rZW4gPSBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIjZmVuZ19zaG93X3Rva2VuIikudmFsdWU7CiAgICAgIGNvbnN0IGRhdGEgPSB7cGhvbmUsIHBhc3N3b3JkLCBjb2RlLCB0aWNrZXQ6IGZlbmdTaG93VG9rZW59OwoKICAgICAgY29uc3QgbWVzc2FnZUVsbSA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3IoIi5yZXN1bHQtd3JhcCAubWVzc2FnZSIpOwogICAgICBtZXNzYWdlRWxtLmlubmVySFRNTCA9ICLmraPlnKjmjojmnYPkuK0sIOivt+eojeWQji4uLiI7CiAgICAgIGNvbnN0IHRva2VuRWxtID0gZG9jdW1lbnQucXVlcnlTZWxlY3RvcigiLnJlc3VsdC13cmFwICNmZW5nX2F1dGhfdG9rZW4iKTsKCiAgICAgIC8vIOWPkei1tyBQT1NUIOivt+axggogICAgICBmZXRjaCgiaHR0cDovL20uZmVuZ3Nob3dzLmNvbS91c2VyL29hdXRoL2xvZ2luIiwgewogICAgICAgIG1ldGhvZDogIlBPU1QiLCAvLyDkvb/nlKggUE9TVCDmlrnms5UKICAgICAgICBoZWFkZXJzOiB7CiAgICAgICAgICAiQ29udGVudC1UeXBlIjogImFwcGxpY2F0aW9uL2pzb24iLCAvLyDmjIflrpror7fmsYLkvZPnsbvlnovkuLogSlNPTgogICAgICAgIH0sCiAgICAgICAgYm9keTogSlNPTi5zdHJpbmdpZnkoZGF0YSksIC8vIOWwhuWvueixoei9rOaNouS4uiBKU09OIOWtl+espuS4sgogICAgICB9KQogICAgICAgIC50aGVuKChyZXNwb25zZSkgPT4gewogICAgICAgICAgaWYgKCFyZXNwb25zZS5vaykgewogICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYEhUVFAgZXJyb3IhIFN0YXR1czogJHtyZXNwb25zZS5zdGF0dXN9YCk7CiAgICAgICAgICB9CiAgICAgICAgICByZXR1cm4gcmVzcG9uc2UuanNvbigpOyAvLyDop6PmnpAgSlNPTiDlk43lupTkvZMKICAgICAgICB9KQogICAgICAgIC50aGVuKChkYXRhKSA9PiB7CiAgICAgICAgICBpZiAoZGF0YS5zdGF0dXMgIT09ICcwJykgewogICAgICAgICAgICB0aHJvdyBuZXcgRXJyb3IoYOaOiOadg+Wksei0pTogJHtkYXRhLm1lc3NhZ2V9YCk7CiAgICAgICAgICB9CiAgICAgICAgICBjb25zdCB0b2tlbiA9IGRhdGE/LmRhdGE/LnRva2VuOwogICAgICAgICAgbWVzc2FnZUVsbS5pbm5lckhUTUwgPSAn5o6I5p2D5oiQ5Yqf77yM6K+35bCG6I635Y+W5Yiw55qE5YC86K6+572u5Yiw56iL5bqP55qE546v5aKD5Y+Y6YeP5Lit77yMa2V5OiBmZW5nX3Rva2VuJzsKICAgICAgICAgIHRva2VuRWxtLnZhbHVlID0gdG9rZW47CiAgICAgICAgfSkKICAgICAgICAuY2F0Y2goKGVycm9yKSA9PiB7CiAgICAgICAgICBtZXNzYWdlRWxtLmlubmVySFRNTCA9IGVycm9yOwogICAgICAgICAgc2V0VGltZW91dCgoKSA9PiBtZXNzYWdlRWxtLmlubmVySFRNTCA9ICfor7floavlhaXoh6rlt7HnmoTlh6Tlh7Dnp4DotKblj7fkv6Hmga/lkI7ngrnlh7vnmbvlvZXojrflj5YgdG9rZW4nLCA1MDAwKTsKICAgICAgICB9KTsKICAgIH07CiAgPC9zY3JpcHQ+CjwvaHRtbD4K`

	// ConfigPageHtml 配置页地址
	ConfigPageHtml = `PCFET0NUWVBFIGh0bWw+CjxodG1sPgogIDxoZWFkPgogICAgPG1ldGEgY2hhcnNldD0iVVRGLTgiIC8+CiAgICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBpbml0aWFsLXNjYWxlPTEuMCIgLz4KICAgIDx0aXRsZT5saXZlLXNlcnZlciDphY3nva7pobU8L3RpdGxlPgogICAgPCEtLSDlvJXlhaUgQm9vdHN0cmFwIENTUyAtLT4KICAgIDxsaW5rCiAgICAgIGhyZWY9Imh0dHBzOi8vY2RuLmpzZGVsaXZyLm5ldC9ucG0vYm9vdHN0cmFwQDUuMy4wL2Rpc3QvY3NzL2Jvb3RzdHJhcC5taW4uY3NzIgogICAgICByZWw9InN0eWxlc2hlZXQiCiAgICAvPgogICAgPHN0eWxlPgogICAgICAucGFnZS1jb250YWluZXIgewogICAgICAgIHBhZGRpbmc6IDIwcHggMTBweDsKICAgICAgfQogICAgICAucGFnZS1jb250YWluZXIgaW5wdXQgewogICAgICAgIG1hcmdpbi1ib3R0b206IDhweDsKICAgICAgfQogICAgPC9zdHlsZT4KICA8L2hlYWQ+CiAgPGJvZHk+CiAgICA8ZGl2IGNsYXNzPSJwYWdlLWNvbnRhaW5lciI+CiAgICAgIDxoMj5saXZlLXNlcnZlciDphY3nva7pobU8L2gyPgoKICAgICAgPCEtLSDlsZXnpLrnlKjmiLfkvKDpgJLnmoTnqIvluo/lr4bpkqUgLS0+CiAgICAgIDxkaXYgY2xhc3M9Im1iLTMgcm93IHNlY3JldC13cmFwIj4KICAgICAgICA8bGFiZWwgY2xhc3M9ImNvbC1zbS0yIGNvbC1mb3JtLWxhYmVsIj7nqIvluo/lr4bpkqXvvJo8L2xhYmVsPgogICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0xMCI+CiAgICAgICAgICA8aW5wdXQgY2xhc3M9ImZvcm0tY29udHJvbCBzZWNyZXQiIHR5cGU9InRleHQiIGRpc2FibGVkIC8+CiAgICAgICAgPC9kaXY+CiAgICAgIDwvZGl2PgoKICAgICAgPCEtLSDnjq/looPlj5jph48gLS0+CiAgICAgIDxkaXYgY2xhc3M9ImVudi13cmFwIj4KICAgICAgICA8aDM+546v5aKD5Y+Y6YePPC9oMz4KICAgICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tNSI+CiAgICAgICAgICAgIDxpbnB1dAogICAgICAgICAgICAgIGNsYXNzPSJmb3JtLWNvbnRyb2wgZW52LWtleSIKICAgICAgICAgICAgICB0eXBlPSJ0ZXh0IgogICAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaUga2V5LCDlpoI6IGZlbmdfdG9rZW4iCiAgICAgICAgICAgIC8+CiAgICAgICAgICA8L2Rpdj4KICAgICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS01Ij4KICAgICAgICAgICAgPGlucHV0CiAgICAgICAgICAgICAgY2xhc3M9ImZvcm0tY29udHJvbCBlbnYtdmFsdWUiCiAgICAgICAgICAgICAgdHlwZT0idGV4dCIKICAgICAgICAgICAgICBwbGFjZWhvbGRlcj0i6K+36L6T5YWlIHZhbHVlIgogICAgICAgICAgICAvPgogICAgICAgICAgPC9kaXY+CiAgICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tMiI+CiAgICAgICAgICAgIDxidXR0b24gdHlwZT0iYnV0dG9uIiBjbGFzcz0iYnRuIGJ0bi1zdWNjZXNzIiBvbmNsaWNrPSJnZXRFbnYoKSI+CiAgICAgICAgICAgICAg5p+lIOivogogICAgICAgICAgICA8L2J1dHRvbj4KICAgICAgICAgICAgPGJ1dHRvbiB0eXBlPSJidXR0b24iIGNsYXNzPSJidG4gYnRuLXByaW1hcnkiIG9uY2xpY2s9InNldEVudigpIj4KICAgICAgICAgICAgICDmt7sg5YqgCiAgICAgICAgICAgIDwvYnV0dG9uPgogICAgICAgICAgPC9kaXY+CiAgICAgICAgPC9kaXY+CiAgICAgIDwvZGl2PgoKICAgICAgPCEtLSDpu5HlkI3ljZUgLS0+CiAgICAgIDxkaXYgY2xhc3M9ImJsYWNrLWlwLXdyYXAiPgogICAgICAgIDxoMz7pu5HlkI3ljZU8L2gzPgogICAgICAgIDxkaXYgY2xhc3M9Im1iLTMgcm93Ij4KICAgICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0xMCI+CiAgICAgICAgICAgIDxpbnB1dAogICAgICAgICAgICAgIGNsYXNzPSJmb3JtLWNvbnRyb2wgYmxhY2staXAiCiAgICAgICAgICAgICAgdHlwZT0idGV4dCIKICAgICAgICAgICAgICBwbGFjZWhvbGRlcj0i6K+36L6T5YWl6KaB5Yqg5YWl6buR5ZCN5Y2V55qEIGlwLCDlpoI6IDEuMS4xLjEiCiAgICAgICAgICAgIC8+CiAgICAgICAgICA8L2Rpdj4KICAgICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS0yIj4KICAgICAgICAgICAgPGJ1dHRvbgogICAgICAgICAgICAgIHR5cGU9ImJ1dHRvbiIKICAgICAgICAgICAgICBjbGFzcz0iYnRuIGJ0bi1wcmltYXJ5IgogICAgICAgICAgICAgIG9uY2xpY2s9InNldEJsYWNrSXAoKSIKICAgICAgICAgICAgPgogICAgICAgICAgICAgIOa3uyDliqAKICAgICAgICAgICAgPC9idXR0b24+CiAgICAgICAgICA8L2Rpdj4KICAgICAgICA8L2Rpdj4KICAgICAgPC9kaXY+CgogICAgICA8IS0tIOWcsOWfn+eZveWQjeWNlSAtLT4KICAgICAgPGRpdiBjbGFzcz0id2hpdGUtYXJlYS13cmFwIj4KICAgICAgICA8aDM+5Zyw5Z+f55m95ZCN5Y2VPC9oMz4KICAgICAgICA8ZGl2IGNsYXNzPSJtYi0zIHJvdyI+CiAgICAgICAgICA8ZGl2IGNsYXNzPSJjb2wtc20tOCI+CiAgICAgICAgICAgIDxpbnB1dAogICAgICAgICAgICAgIGNsYXNzPSJmb3JtLWNvbnRyb2wgd2hpdGUtYXJlYSIKICAgICAgICAgICAgICB0eXBlPSJ0ZXh0IgogICAgICAgICAgICAgIHBsYWNlaG9sZGVyPSLor7fovpPlhaXopoHliqDlhaXnmb3lkI3ljZXnmoTlnLDln58sIOWmgjog5bm/5LicL+S9m+WxsS/ljZfmtbciCiAgICAgICAgICAgIC8+CiAgICAgICAgICA8L2Rpdj4KICAgICAgICAgIDxkaXYgY2xhc3M9ImNvbC1zbS00Ij4KICAgICAgICAgICAgPGJ1dHRvbgogICAgICAgICAgICAgIHR5cGU9ImJ1dHRvbiIKICAgICAgICAgICAgICBjbGFzcz0iYnRuIGJ0bi1wcmltYXJ5IgogICAgICAgICAgICAgIG9uY2xpY2s9InNldFdoaXRlQXJlYSgnc2V0JykiCiAgICAgICAgICAgID4KICAgICAgICAgICAgICDmt7sg5YqgCiAgICAgICAgICAgIDwvYnV0dG9uPgogICAgICAgICAgICA8YnV0dG9uCiAgICAgICAgICAgICAgdHlwZT0iYnV0dG9uIgogICAgICAgICAgICAgIGNsYXNzPSJidG4gYnRuLWRhbmdlciIKICAgICAgICAgICAgICBvbmNsaWNrPSJzZXRXaGl0ZUFyZWEoJ2RlbCcpIgogICAgICAgICAgICA+CiAgICAgICAgICAgICAg56e7IOmZpAogICAgICAgICAgICA8L2J1dHRvbj4KICAgICAgICAgIDwvZGl2PgogICAgICAgIDwvZGl2PgogICAgICA8L2Rpdj4KICAgIDwvZGl2PgogIDwvYm9keT4KCiAgPHNjcmlwdD4KICAgIC8vIOiOt+WPluWfuuacrOWPguaVsAogICAgY29uc3QgZ2V0QmFzZVBhcmFtcyA9ICgpID0+IHsKICAgICAgY29uc3Qgc2VjcmV0ID0gbmV3IFVSTCh3aW5kb3cubG9jYXRpb24uaHJlZikuc2VhcmNoUGFyYW1zLmdldCgic2VjcmV0Iik7CiAgICAgIHJldHVybiB7CiAgICAgICAgc2VjcmV0LCAvLyDnqIvluo/lr4bpkqUKICAgICAgICBob3N0OiAiIiwgLy8g5Y+R6YCB6K+35rGC55qE5Li75py65Zyw5Z2ACiAgICAgIH07CiAgICB9OwoKICAgIHdpbmRvdy5vbmxvYWQgPSAoKSA9PiB7CiAgICAgIC8vIOiuvue9rueoi+W6j+WvhumSpQogICAgICBjb25zdCB7IHNlY3JldCB9ID0gZ2V0QmFzZVBhcmFtcygpOwogICAgICBkb2N1bWVudC5xdWVyeVNlbGVjdG9yKCIucGFnZS1jb250YWluZXIgLnNlY3JldCIpLnZhbHVlID0gc2VjcmV0OwogICAgfTsKCiAgICAvLyDojrflj5bovpPlhaXmoYbnmoTlgLwKICAgIGNvbnN0IGdldElucHV0VmFsdWUgPSAoc2VsZWN0b3IgPSAiIikgPT4gewogICAgICByZXR1cm4gKGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3Ioc2VsZWN0b3IpLnZhbHVlIHx8ICIiKS50cmltKCk7CiAgICB9OwoKICAgIC8vIOiuvue9rui+k+WFpeahhueahOWAvAogICAgY29uc3Qgc2V0SW5wdXRWYWx1ZSA9IChzZWxlY3RvciA9ICIiLCB2YWx1ZSA9ICIiKSA9PiB7CiAgICAgIGNvbnN0IHRhcmdldCA9IGRvY3VtZW50LnF1ZXJ5U2VsZWN0b3Ioc2VsZWN0b3IpOwogICAgICBpZiAoIXRhcmdldCkgewogICAgICAgIHJldHVybgogICAgICB9CiAgICAgIHRhcmdldC52YWx1ZSA9IHZhbHVlOwogICAgfQoKICAgIC8vIEdFVCDor7fmsYLmjIflrprlnLDlnYAKICAgIGNvbnN0IGZldGNoVXJsID0gKAogICAgICB1cmwgPSAiIiwKICAgICAgc3VjY2Vzc0Z1bmMgPSAoKSA9PiB7fSwKICAgICAgZXJyb3JGdW5jID0gKGVycikgPT4ge30KICAgICkgPT4gewogICAgICBmZXRjaCh1cmwpCiAgICAgICAgLnRoZW4oKHJlcykgPT4gewogICAgICAgICAgaWYgKHJlcy5vaykgewogICAgICAgICAgICByZXR1cm4gc3VjY2Vzc0Z1bmMoKTsKICAgICAgICAgIH0KICAgICAgICAgIGVycm9yRnVuYyhyZXMuc3RhdHVzVGV4dCk7CiAgICAgICAgfSkKICAgICAgICAuY2F0Y2goKGVycikgPT4gewogICAgICAgICAgZXJyb3JGdW5jKGVycik7CiAgICAgICAgfSk7CiAgICB9OwoKICAgIC8vIOafpeivoueOr+Wig+WPmOmHjwogICAgY29uc3QgZ2V0RW52ID0gKCkgPT4gewogICAgICBjb25zdCBlbnZLZXkgPSBnZXRJbnB1dFZhbHVlKCIuZW52LXdyYXAgLmVudi1rZXkiKTsKICAgICAgaWYgKGVudktleSA9PT0gIiIpIHsKICAgICAgICByZXR1cm4gYWxlcnQoIui+k+WFpeWPguaVsOS4jeiDveS4uuepuiIpOwogICAgICB9CgogICAgICBjb25zdCB7IHNlY3JldCwgaG9zdCB9ID0gZ2V0QmFzZVBhcmFtcygpOwogICAgICBjb25zdCB1cmwgPSBgJHtob3N0fS9lbnY/a2V5PSR7ZW52S2V5fSZzZWNyZXQ9JHtzZWNyZXR9YDsKICAgICAgZmV0Y2godXJsKS50aGVuKHJlcyA9PiB7CiAgICAgICAgaWYgKCFyZXMub2spIHsKICAgICAgICAgIHJldHVybiBhbGVydChg5p+l6K+i5aSx6LSlOiAke3Jlcy5zdGF0dXNUZXh0fWApOwogICAgICAgIH0KICAgICAgICByZXR1cm4gcmVzLnRleHQoKTsKICAgICAgfSkudGhlbihyZXNwVGV4dCA9PiB7CiAgICAgICAgc2V0SW5wdXRWYWx1ZSgiLmVudi13cmFwIC5lbnYtdmFsdWUiLCByZXNwVGV4dCk7CiAgICAgIH0pLmNhdGNoKGVyciA9PiB7CiAgICAgICAgYWxlcnQoYOafpeivouWksei0pTogJHtlcnJ9YCk7CiAgICAgIH0pCiAgICB9CgogICAgLy8g6K6+572u546v5aKD5Y+Y6YePCiAgICBjb25zdCBzZXRFbnYgPSAoKSA9PiB7CiAgICAgIGNvbnN0IGVudktleSA9IGdldElucHV0VmFsdWUoIi5lbnYtd3JhcCAuZW52LWtleSIpOwogICAgICBjb25zdCBlbnZWYWx1ZSA9IGdldElucHV0VmFsdWUoIi5lbnYtd3JhcCAuZW52LXZhbHVlIik7CiAgICAgIGlmIChlbnZLZXkgPT09ICIiIHx8IGVudlZhbHVlID09PSAiIikgewogICAgICAgIHJldHVybiBhbGVydCgi6L6T5YWl5Y+C5pWw5LiN6IO95Li656m6Iik7CiAgICAgIH0KCiAgICAgIGNvbnN0IHsgc2VjcmV0LCBob3N0IH0gPSBnZXRCYXNlUGFyYW1zKCk7CiAgICAgIGNvbnN0IHVybCA9IGAke2hvc3R9L2Vudj9zZWNyZXQ9JHtzZWNyZXR9YDsKICAgICAgY29uc3QgZm9ybURhdGEgPSBuZXcgRm9ybURhdGEoKTsKICAgICAgZm9ybURhdGEuYXBwZW5kKCJrZXkiLCBlbnZLZXkpOwogICAgICBmb3JtRGF0YS5hcHBlbmQoInZhbHVlIiwgZW52VmFsdWUpOwoKICAgICAgZmV0Y2godXJsLCB7CiAgICAgICAgbWV0aG9kOiAiUE9TVCIsCiAgICAgICAgYm9keTogZm9ybURhdGEsCiAgICAgIH0pCiAgICAgICAgLnRoZW4oKHJlcykgPT4gewogICAgICAgICAgaWYgKHJlcy5vaykgewogICAgICAgICAgICByZXR1cm4gYWxlcnQoIua3u+WKoOaIkOWKnyIpOwogICAgICAgICAgfQogICAgICAgICAgYWxlcnQoYOa3u+WKoOWksei0pTogJHtyZXMuc3RhdHVzVGV4dH1gKTsKICAgICAgICB9KQogICAgICAgIC5jYXRjaCgoZXJyKSA9PiB7CiAgICAgICAgICBhbGVydChg5re75Yqg5aSx6LSlOiAke2Vycn1gKTsKICAgICAgICB9KTsKICAgIH07CgogICAgLy8g6K6+572u6buR5ZCN5Y2VIGlwCiAgICBjb25zdCBzZXRCbGFja0lwID0gKCkgPT4gewogICAgICBjb25zdCBibGFja0lwID0gZ2V0SW5wdXRWYWx1ZSgiLmJsYWNrLWlwLXdyYXAgLmJsYWNrLWlwIik7CiAgICAgIGlmIChibGFja0lwID09PSAiIikgewogICAgICAgIHJldHVybiBhbGVydCgi5Y+C5pWw5LiN6IO95Li656m6Iik7CiAgICAgIH0KICAgICAgY29uc3QgeyBzZWNyZXQsIGhvc3QgfSA9IGdldEJhc2VQYXJhbXMoKTsKICAgICAgY29uc3QgdXJsID0gYCR7aG9zdH0vYmxhY2tfaXA/aXA9JHtibGFja0lwfSZzZWNyZXQ9JHtzZWNyZXR9YDsKICAgICAgZmV0Y2hVcmwoCiAgICAgICAgdXJsLAogICAgICAgICgpID0+IGFsZXJ0KCLmt7vliqDmiJDlip8iKSwKICAgICAgICAoZXJyKSA9PiBhbGVydChg5re75Yqg5aSx6LSlOiAke2Vycn1gKQogICAgICApOwogICAgfTsKCiAgICAvLyDorr7nva7lnLDln5/nmb3lkI3ljZUKICAgIGNvbnN0IHNldFdoaXRlQXJlYSA9IChhY3Rpb24gPSAic2V0IikgPT4gewogICAgICBjb25zdCB3aGl0ZUFyZWEgPSBnZXRJbnB1dFZhbHVlKCIud2hpdGUtYXJlYS13cmFwIC53aGl0ZS1hcmVhIik7CiAgICAgIGlmICh3aGl0ZUFyZWEgPT09ICIiKSB7CiAgICAgICAgcmV0dXJuIGFsZXJ0KCLovpPlhaXlj4LmlbDkuI3og73kuLrnqboiKTsKICAgICAgfQogICAgICBjb25zdCBhY3Rpb25UaXAgPSBhY3Rpb24gPT09ICJzZXQiID8gIua3u+WKoCIgOiAi56e76ZmkIjsKICAgICAgY29uc3QgeyBzZWNyZXQsIGhvc3QgfSA9IGdldEJhc2VQYXJhbXMoKTsKICAgICAgY29uc3QgdXJsID0gYCR7aG9zdH0vd2hpdGVfYXJlYS8ke2FjdGlvbn0/YXJlYT0ke3doaXRlQXJlYX0mc2VjcmV0PSR7c2VjcmV0fWA7CiAgICAgIGZldGNoVXJsKAogICAgICAgIHVybCwKICAgICAgICAoKSA9PiBhbGVydChgJHthY3Rpb25UaXB95oiQ5YqfYCksCiAgICAgICAgKGVycikgPT4gYWxlcnQoYCR7YWN0aW9uVGlwfeWksei0pTogJHtlcnJ9YCkKICAgICAgKTsKICAgIH07CiAgPC9zY3JpcHQ+CjwvaHRtbD4K`
)

const (
	Gin_IpAddrInfoKey = "ip_addr_info" // 在 gin 上下文中存放 IP 地址信息的键
)
