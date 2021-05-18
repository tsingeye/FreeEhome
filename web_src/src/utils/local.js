export const getLocal = (key) => {
    let value = localStorage.getItem(key) || '';
    if (value.startsWith('[') || value.includes('{')) {
        return JSON.parse(value);
    } else {
        return value
    }
}
export const setLocal = (key,value)=>{
    if(typeof value == 'object'){
        value = JSON.stringify(value);
    }
    localStorage.setItem(key,value);
}