const fs = require('fs');
const file = 'c:/Users/Admin/Documents/GitHub/THD_iOS_MDM_Portal/frontend/src/features/device-profiles/components/ProfilesList.tsx';
let content = fs.readFileSync(file, 'utf8');

// Replace orange
content = content.replace(/className=\"flex items-center gap-3 w-full p-2\.5 hover:bg-orange-50 rounded-lg text-slate-700 hover:text-orange-700 transition-all duration-200 text-sm text-left font-medium hover:font-semibold cursor-pointer hover:shadow-sm\"\n\s*>\n\s*<([a-zA-Z0-9]+) className=\"w-4 h-4 text-slate-400\" \/>/g, 'className=\"group flex items-center gap-3 w-full p-2.5 hover:bg-orange-100 rounded-lg text-slate-700 hover:text-orange-700 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-sm\"\n                                            >\n                                                <$1 className=\"w-4 h-4 text-slate-400 group-hover:text-orange-600 transition-colors\" />');

// Replace slate
content = content.replace(/className=\"flex items-center gap-3 w-full p-2\.5 hover:bg-slate-100 rounded-lg text-slate-700 hover:text-slate-900 transition-all duration-200 text-sm text-left font-medium hover:font-semibold cursor-pointer hover:shadow-sm\"\n\s*>\n\s*<([a-zA-Z0-9]+) className=\"w-4 h-4 text-slate-400\" \/>/g, 'className=\"group flex items-center gap-3 w-full p-2.5 hover:bg-slate-200 rounded-lg text-slate-700 hover:text-slate-900 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-sm\"\n                                            >\n                                                <$1 className=\"w-4 h-4 text-slate-400 group-hover:text-slate-700 transition-colors\" />');

// Replace rose
content = content.replace(/className=\"flex items-center gap-3 w-full p-2\.5 hover:bg-rose-50 rounded-lg text-slate-700 hover:text-rose-700 transition-all duration-200 text-sm text-left font-medium hover:font-semibold cursor-pointer hover:shadow-sm\"\n\s*>\n\s*<([a-zA-Z0-9]+) className=\"w-4 h-4 text-slate-400\" \/>/g, 'className=\"group flex items-center gap-3 w-full p-2.5 hover:bg-rose-100 rounded-lg text-slate-700 hover:text-rose-700 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-sm\"\n                                            >\n                                                <$1 className=\"w-4 h-4 text-slate-400 group-hover:text-rose-600 transition-colors\" />');

// Replace teal
content = content.replace(/className=\"flex items-center gap-3 w-full p-2\.5 hover:bg-teal-50 rounded-lg text-slate-700 hover:text-teal-700 transition-all duration-200 text-sm text-left font-medium hover:font-semibold cursor-pointer hover:shadow-sm\"\n\s*>\n\s*<([a-zA-Z0-9]+) className=\"w-4 h-4 text-slate-400\" \/>/g, 'className=\"group flex items-center gap-3 w-full p-2.5 hover:bg-teal-100 rounded-lg text-slate-700 hover:text-teal-700 transition-all duration-200 text-sm text-left font-medium hover:font-bold cursor-pointer hover:shadow-sm\"\n                                            >\n                                                <$1 className=\"w-4 h-4 text-slate-400 group-hover:text-teal-600 transition-colors\" />');

fs.writeFileSync(file, content);
