
export declare interface StdResponse {
    code: number
    msg:string
    data:any
}


export const NotErrCode = function (res:StdResponse):boolean{
    if (res.code == 200 ){
        return true
    }
    console.error(`error: code:${res.code},msg: ${res.msg}`)
    return false
}




export const Panic = function (error:any){
    console.error("error: ",error)
}