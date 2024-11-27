import{c as H,e as _}from"./chunk-VHNPENGG.js";import{b as y,c as G,f as d,g as s,h as O,i as l,j as T,k as Q,l as L}from"./chunk-33KK2FKQ.js";import{a as N}from"./chunk-HB55W55I.js";import{Hb as F,Sb as q,a as S,b as C,h as E,pc as I,qa as W,ua as v}from"./chunk-FJILXII2.js";var Yt=(()=>{class r{constructor(){this.themeInfo=W(_),this.transloco=W(N),this.$data=new E,this.width=500,this.height=500}ngOnInit(){this.updateChart(),this.$data.subscribe(e=>{this.data=e,this.updateChart()}),this.themeInfo.info$.subscribe(()=>{this.updateChart()}),this.transloco.langChanges$.subscribe(()=>{this.updateChart()})}updateChart(){this.chartConfig=this.adapter.create(this.data)}static{this.\u0275fac=function(n){return new(n||r)}}static{this.\u0275cmp=v({type:r,selectors:[["app-chart"]],inputs:{$data:"$data",adapter:"adapter",width:"width",height:"height"},standalone:!0,features:[I],decls:1,vars:5,consts:[["baseChart","",3,"data","options","type","height","width"]],template:function(n,a){n&1&&q(0,"canvas",0),n&2&&F("data",a.chartConfig.data)("options",a.chartConfig.options)("type",a.chartConfig.type)("height",a.height)("width",a.width)},dependencies:[H]})}}return r})();function M(r,t){let e=s(r,t?.in);return e.setHours(0,0,0,0),e}function $(r,t,e){let[n,a]=Q(e?.in,r,t),o=M(n),c=M(a),m=+o-T(o),h=+c-T(c);return Math.round((m-h)/G)}function B(r,t){let e=s(r,t?.in);return e.setFullYear(e.getFullYear(),0,1),e.setHours(0,0,0,0),e}function X(r,t){let e=s(r,t?.in);return $(e,B(e))+1}function p(r,t){return l(r,C(S({},t),{weekStartsOn:1}))}function b(r,t){let e=s(r,t?.in),n=e.getFullYear(),a=d(e,0);a.setFullYear(n+1,0,4),a.setHours(0,0,0,0);let o=p(a),c=d(e,0);c.setFullYear(n,0,4),c.setHours(0,0,0,0);let m=p(c);return e.getTime()>=o.getTime()?n+1:e.getTime()>=m.getTime()?n:n-1}function R(r,t){let e=b(r,t),n=d(t?.in||r,0);return n.setFullYear(e,0,4),n.setHours(0,0,0,0),p(n)}function j(r,t){let e=s(r,t?.in),n=+p(e)-+R(e);return Math.round(n/y)+1}function k(r,t){let e=s(r,t?.in),n=e.getFullYear(),a=O(),o=t?.firstWeekContainsDate??t?.locale?.options?.firstWeekContainsDate??a.firstWeekContainsDate??a.locale?.options?.firstWeekContainsDate??1,c=d(t?.in||r,0);c.setFullYear(n+1,0,o),c.setHours(0,0,0,0);let m=l(c,t),h=d(t?.in||r,0);h.setFullYear(n,0,o),h.setHours(0,0,0,0);let D=l(h,t);return+e>=+m?n+1:+e>=+D?n:n-1}function V(r,t){let e=O(),n=t?.firstWeekContainsDate??t?.locale?.options?.firstWeekContainsDate??e.firstWeekContainsDate??e.locale?.options?.firstWeekContainsDate??1,a=k(r,t),o=d(t?.in||r,0);return o.setFullYear(a,0,n),o.setHours(0,0,0,0),l(o,t)}function Z(r,t){let e=s(r,t?.in),n=+l(e,t)-+V(e,t);return Math.round(n/y)+1}function i(r,t){let e=r<0?"-":"",n=Math.abs(r).toString().padStart(t,"0");return e+n}var g={y(r,t){let e=r.getFullYear(),n=e>0?e:1-e;return i(t==="yy"?n%100:n,t.length)},M(r,t){let e=r.getMonth();return t==="M"?String(e+1):i(e+1,2)},d(r,t){return i(r.getDate(),t.length)},a(r,t){let e=r.getHours()/12>=1?"pm":"am";switch(t){case"a":case"aa":return e.toUpperCase();case"aaa":return e;case"aaaaa":return e[0];case"aaaa":default:return e==="am"?"a.m.":"p.m."}},h(r,t){return i(r.getHours()%12||12,t.length)},H(r,t){return i(r.getHours(),t.length)},m(r,t){return i(r.getMinutes(),t.length)},s(r,t){return i(r.getSeconds(),t.length)},S(r,t){let e=t.length,n=r.getMilliseconds(),a=Math.trunc(n*Math.pow(10,e-3));return i(a,t.length)}};var x={am:"am",pm:"pm",midnight:"midnight",noon:"noon",morning:"morning",afternoon:"afternoon",evening:"evening",night:"night"},P={G:function(r,t,e){let n=r.getFullYear()>0?1:0;switch(t){case"G":case"GG":case"GGG":return e.era(n,{width:"abbreviated"});case"GGGGG":return e.era(n,{width:"narrow"});case"GGGG":default:return e.era(n,{width:"wide"})}},y:function(r,t,e){if(t==="yo"){let n=r.getFullYear(),a=n>0?n:1-n;return e.ordinalNumber(a,{unit:"year"})}return g.y(r,t)},Y:function(r,t,e,n){let a=k(r,n),o=a>0?a:1-a;if(t==="YY"){let c=o%100;return i(c,2)}return t==="Yo"?e.ordinalNumber(o,{unit:"year"}):i(o,t.length)},R:function(r,t){let e=b(r);return i(e,t.length)},u:function(r,t){let e=r.getFullYear();return i(e,t.length)},Q:function(r,t,e){let n=Math.ceil((r.getMonth()+1)/3);switch(t){case"Q":return String(n);case"QQ":return i(n,2);case"Qo":return e.ordinalNumber(n,{unit:"quarter"});case"QQQ":return e.quarter(n,{width:"abbreviated",context:"formatting"});case"QQQQQ":return e.quarter(n,{width:"narrow",context:"formatting"});case"QQQQ":default:return e.quarter(n,{width:"wide",context:"formatting"})}},q:function(r,t,e){let n=Math.ceil((r.getMonth()+1)/3);switch(t){case"q":return String(n);case"qq":return i(n,2);case"qo":return e.ordinalNumber(n,{unit:"quarter"});case"qqq":return e.quarter(n,{width:"abbreviated",context:"standalone"});case"qqqqq":return e.quarter(n,{width:"narrow",context:"standalone"});case"qqqq":default:return e.quarter(n,{width:"wide",context:"standalone"})}},M:function(r,t,e){let n=r.getMonth();switch(t){case"M":case"MM":return g.M(r,t);case"Mo":return e.ordinalNumber(n+1,{unit:"month"});case"MMM":return e.month(n,{width:"abbreviated",context:"formatting"});case"MMMMM":return e.month(n,{width:"narrow",context:"formatting"});case"MMMM":default:return e.month(n,{width:"wide",context:"formatting"})}},L:function(r,t,e){let n=r.getMonth();switch(t){case"L":return String(n+1);case"LL":return i(n+1,2);case"Lo":return e.ordinalNumber(n+1,{unit:"month"});case"LLL":return e.month(n,{width:"abbreviated",context:"standalone"});case"LLLLL":return e.month(n,{width:"narrow",context:"standalone"});case"LLLL":default:return e.month(n,{width:"wide",context:"standalone"})}},w:function(r,t,e,n){let a=Z(r,n);return t==="wo"?e.ordinalNumber(a,{unit:"week"}):i(a,t.length)},I:function(r,t,e){let n=j(r);return t==="Io"?e.ordinalNumber(n,{unit:"week"}):i(n,t.length)},d:function(r,t,e){return t==="do"?e.ordinalNumber(r.getDate(),{unit:"date"}):g.d(r,t)},D:function(r,t,e){let n=X(r);return t==="Do"?e.ordinalNumber(n,{unit:"dayOfYear"}):i(n,t.length)},E:function(r,t,e){let n=r.getDay();switch(t){case"E":case"EE":case"EEE":return e.day(n,{width:"abbreviated",context:"formatting"});case"EEEEE":return e.day(n,{width:"narrow",context:"formatting"});case"EEEEEE":return e.day(n,{width:"short",context:"formatting"});case"EEEE":default:return e.day(n,{width:"wide",context:"formatting"})}},e:function(r,t,e,n){let a=r.getDay(),o=(a-n.weekStartsOn+8)%7||7;switch(t){case"e":return String(o);case"ee":return i(o,2);case"eo":return e.ordinalNumber(o,{unit:"day"});case"eee":return e.day(a,{width:"abbreviated",context:"formatting"});case"eeeee":return e.day(a,{width:"narrow",context:"formatting"});case"eeeeee":return e.day(a,{width:"short",context:"formatting"});case"eeee":default:return e.day(a,{width:"wide",context:"formatting"})}},c:function(r,t,e,n){let a=r.getDay(),o=(a-n.weekStartsOn+8)%7||7;switch(t){case"c":return String(o);case"cc":return i(o,t.length);case"co":return e.ordinalNumber(o,{unit:"day"});case"ccc":return e.day(a,{width:"abbreviated",context:"standalone"});case"ccccc":return e.day(a,{width:"narrow",context:"standalone"});case"cccccc":return e.day(a,{width:"short",context:"standalone"});case"cccc":default:return e.day(a,{width:"wide",context:"standalone"})}},i:function(r,t,e){let n=r.getDay(),a=n===0?7:n;switch(t){case"i":return String(a);case"ii":return i(a,t.length);case"io":return e.ordinalNumber(a,{unit:"day"});case"iii":return e.day(n,{width:"abbreviated",context:"formatting"});case"iiiii":return e.day(n,{width:"narrow",context:"formatting"});case"iiiiii":return e.day(n,{width:"short",context:"formatting"});case"iiii":default:return e.day(n,{width:"wide",context:"formatting"})}},a:function(r,t,e){let a=r.getHours()/12>=1?"pm":"am";switch(t){case"a":case"aa":return e.dayPeriod(a,{width:"abbreviated",context:"formatting"});case"aaa":return e.dayPeriod(a,{width:"abbreviated",context:"formatting"}).toLowerCase();case"aaaaa":return e.dayPeriod(a,{width:"narrow",context:"formatting"});case"aaaa":default:return e.dayPeriod(a,{width:"wide",context:"formatting"})}},b:function(r,t,e){let n=r.getHours(),a;switch(n===12?a=x.noon:n===0?a=x.midnight:a=n/12>=1?"pm":"am",t){case"b":case"bb":return e.dayPeriod(a,{width:"abbreviated",context:"formatting"});case"bbb":return e.dayPeriod(a,{width:"abbreviated",context:"formatting"}).toLowerCase();case"bbbbb":return e.dayPeriod(a,{width:"narrow",context:"formatting"});case"bbbb":default:return e.dayPeriod(a,{width:"wide",context:"formatting"})}},B:function(r,t,e){let n=r.getHours(),a;switch(n>=17?a=x.evening:n>=12?a=x.afternoon:n>=4?a=x.morning:a=x.night,t){case"B":case"BB":case"BBB":return e.dayPeriod(a,{width:"abbreviated",context:"formatting"});case"BBBBB":return e.dayPeriod(a,{width:"narrow",context:"formatting"});case"BBBB":default:return e.dayPeriod(a,{width:"wide",context:"formatting"})}},h:function(r,t,e){if(t==="ho"){let n=r.getHours()%12;return n===0&&(n=12),e.ordinalNumber(n,{unit:"hour"})}return g.h(r,t)},H:function(r,t,e){return t==="Ho"?e.ordinalNumber(r.getHours(),{unit:"hour"}):g.H(r,t)},K:function(r,t,e){let n=r.getHours()%12;return t==="Ko"?e.ordinalNumber(n,{unit:"hour"}):i(n,t.length)},k:function(r,t,e){let n=r.getHours();return n===0&&(n=24),t==="ko"?e.ordinalNumber(n,{unit:"hour"}):i(n,t.length)},m:function(r,t,e){return t==="mo"?e.ordinalNumber(r.getMinutes(),{unit:"minute"}):g.m(r,t)},s:function(r,t,e){return t==="so"?e.ordinalNumber(r.getSeconds(),{unit:"second"}):g.s(r,t)},S:function(r,t){return g.S(r,t)},X:function(r,t,e){let n=r.getTimezoneOffset();if(n===0)return"Z";switch(t){case"X":return J(n);case"XXXX":case"XX":return w(n);case"XXXXX":case"XXX":default:return w(n,":")}},x:function(r,t,e){let n=r.getTimezoneOffset();switch(t){case"x":return J(n);case"xxxx":case"xx":return w(n);case"xxxxx":case"xxx":default:return w(n,":")}},O:function(r,t,e){let n=r.getTimezoneOffset();switch(t){case"O":case"OO":case"OOO":return"GMT"+A(n,":");case"OOOO":default:return"GMT"+w(n,":")}},z:function(r,t,e){let n=r.getTimezoneOffset();switch(t){case"z":case"zz":case"zzz":return"GMT"+A(n,":");case"zzzz":default:return"GMT"+w(n,":")}},t:function(r,t,e){let n=Math.trunc(+r/1e3);return i(n,t.length)},T:function(r,t,e){return i(+r,t.length)}};function A(r,t=""){let e=r>0?"-":"+",n=Math.abs(r),a=Math.trunc(n/60),o=n%60;return o===0?e+String(a):e+String(a)+t+i(o,2)}function J(r,t){return r%60===0?(r>0?"-":"+")+i(Math.abs(r)/60,2):w(r,t)}function w(r,t=""){let e=r>0?"-":"+",n=Math.abs(r),a=i(Math.trunc(n/60),2),o=i(n%60,2);return e+a+t+o}var K=(r,t)=>{switch(r){case"P":return t.date({width:"short"});case"PP":return t.date({width:"medium"});case"PPP":return t.date({width:"long"});case"PPPP":default:return t.date({width:"full"})}},U=(r,t)=>{switch(r){case"p":return t.time({width:"short"});case"pp":return t.time({width:"medium"});case"ppp":return t.time({width:"long"});case"pppp":default:return t.time({width:"full"})}},ot=(r,t)=>{let e=r.match(/(P+)(p+)?/)||[],n=e[1],a=e[2];if(!a)return K(r,t);let o;switch(n){case"P":o=t.dateTime({width:"short"});break;case"PP":o=t.dateTime({width:"medium"});break;case"PPP":o=t.dateTime({width:"long"});break;case"PPPP":default:o=t.dateTime({width:"full"});break}return o.replace("{{date}}",K(n,t)).replace("{{time}}",U(a,t))},z={p:U,P:ot};var it=/^D+$/,st=/^Y+$/,ct=["D","DD","YY","YYYY"];function tt(r){return it.test(r)}function et(r){return st.test(r)}function rt(r,t,e){let n=ut(r,t,e);if(console.warn(n),ct.includes(r))throw new RangeError(n)}function ut(r,t,e){let n=r[0]==="Y"?"years":"days of the month";return`Use \`${r.toLowerCase()}\` instead of \`${r}\` (in \`${t}\`) for formatting ${n} to the input \`${e}\`; see: https://github.com/date-fns/date-fns/blob/master/docs/unicodeTokens.md`}function nt(r){return r instanceof Date||typeof r=="object"&&Object.prototype.toString.call(r)==="[object Date]"}function at(r){return!(!nt(r)&&typeof r!="number"||isNaN(+s(r)))}var ft=/[yYQqMLwIdDecihHKkms]o|(\w)\1*|''|'(''|[^'])+('|$)|./g,dt=/P+p+|P+|p+|''|'(''|[^'])+('|$)|./g,mt=/^'([^]*?)'?$/,ht=/''/g,lt=/[a-zA-Z]/;function Le(r,t,e){let n=O(),a=e?.locale??n.locale??L,o=e?.firstWeekContainsDate??e?.locale?.options?.firstWeekContainsDate??n.firstWeekContainsDate??n.locale?.options?.firstWeekContainsDate??1,c=e?.weekStartsOn??e?.locale?.options?.weekStartsOn??n.weekStartsOn??n.locale?.options?.weekStartsOn??0,m=s(r,e?.in);if(!at(m))throw new RangeError("Invalid time value");let h=t.match(dt).map(f=>{let u=f[0];if(u==="p"||u==="P"){let Y=z[u];return Y(f,a.formatLong)}return f}).join("").match(ft).map(f=>{if(f==="''")return{isToken:!1,value:"'"};let u=f[0];if(u==="'")return{isToken:!1,value:gt(f)};if(P[u])return{isToken:!0,value:f};if(u.match(lt))throw new RangeError("Format string contains an unescaped latin alphabet character `"+u+"`");return{isToken:!1,value:f}});a.localize.preprocessor&&(h=a.localize.preprocessor(m,h));let D={firstWeekContainsDate:o,weekStartsOn:c,locale:a};return h.map(f=>{if(!f.isToken)return f.value;let u=f.value;(!e?.useAdditionalWeekYearTokens&&et(u)||!e?.useAdditionalDayOfYearTokens&&tt(u))&&rt(u,t,String(r));let Y=P[u[0]];return Y(m,u,a.localize,D)}).join("")}function gt(r){let t=r.match(mt);return t?t[1].replace(ht,"'"):r}var $e=(r,t)=>`${r}-${t}`;export{Le as a,Yt as b,$e as c};
