$(function(){
    $('#powerPanel').on('change', function(){
        console.log('changed', $(this).is(':checked'))
        send('/power/0', {on: $(this).is(':checked')})
    })

    $('#btn-cct-group').on('click', function(){
        send('/cct/0', {cct: $('#cctGroup').val()})
    })
    $('#btn-dim-group').on('click', function(){
        send('/dim/0', {dim: $('#dimGroup').val()})
    })

    //operations for light1
    $('#powerLight1').on('change', function(){
        console.log('light1 power changed to', $(this).is(':checked'))
        send('/power/1', {on: $(this).is(':checked')})
    })
    $('#btn-dim-light1').on('click', function(){
        send('/dim/1', {dim: $('#dimLight1').val()})
    })
    $('#btn-cct-light1').on('click', function(){
        send('/cct/1', {cct: $('#cctLight1').val()})
    })
    //operations for light2   
    $('#powerLight2').on('change', function(){
        console.log('light2 power changed to', $(this).is(':checked'))
        send('/power/2', {on: $(this).is(':checked')})
    })
    $('#btn-dim-light2').on('click', function(){
        send('/dim/2', {dim: $('#dimLight2').val()})
    })
    $('#btn-cct-light2').on('click', function(){
        send('/cct/2', {cct: $('#cctLight2').val()})
    })

    $('#statusLight1').on('click', function(){
        updateLight(ls1)
    })
    $('#statusLight2').on('click', function(){
        updateLight(ls2)
    })

    Vue.component('status-light', {
        template: sStatusLight
    })
    new Vue({
        el: '#statusLight1'
    })
    new Vue({
        el: '#statusLight2'
    })

    getStatus()    
})

//<h3>DIM  <span id="statusDim" style="margin-left: 13em;">100</span>%</h3>
var sStatusLight = `
<div>
    <h5><i>Last updated at <span id='statusTime'></span></i></h5>
    <table class="table">
        <tbody>
            <tr>
                <td><h4>INPUT POWER</h4></td>
                <td><span id="statusInputPower">100</span>W</td>
            </tr>
            <tr>
                <td><h4>DIM</h4></td>
                <td><span id="statusDim">100</span>%</td>
            </tr>
            <tr>
                <td><h4>OUTPUT VOLTAGE</h4></td>
                <td><span id="statusOutputVoltage">50</span>V</td>
            </tr>
            <tr>
                <td><h4>OUTPUT CURRENT</h4></td>
                <td><span id="statusOutputCurrent">1.4</span>A</td>
            </tr>
            <tr>
                <td><h4>INTERNAL TEMPERATURE</h4></td>
                <td><span id="statusLoracardTemperature">16</span>Cels</td>
            </tr>
        </tbody>
    </table>
</div>`

var ls1 = null, ls2 = null // light status
function getStatus(){
    var ws = new WebSocket('ws://'+location.host+'/light/status')

    ws.onopen = function(e){
        console.log('connected')
        ws.send('connected')
    }
    ws.onclose = function(e){
        console.log('connetion closed')
    }
    ws.onerror = function(e){
        console.log('connetion error')
    }
    ws.onmessage = function(e){
        data = JSON.parse(e.data)
        console.log('data received', data)

        if(data.id == 1){
            ls1 = data
        }else  if(data.id == 2){
            ls2 = data
        }
    }
}

function updateLight(data){
    if (data == null) return

    $('#statusInputPower').val(data.inputPower)
    $('#statusDim').val(data.dim)
    $('#statusOuptputVoltage').val(data.outputVoltage)
    $('#statusOutputCurrent').val(data.outputCurrent)
    $('#statusTime').val(data.timestamp)
}

function send(url, data){
    console.log(url, data)
    $.ajax({
        type: 'POST',
        url: url,
        data: data,
        success: function(resp){
            swal('Done', getDescName(data)+' operation done successfully', 'success')
        }
    })
}

function getDescName(data){
    if(data.hasOwnProperty('on')) return 'power'
    if(data.hasOwnProperty('dim')) return 'dim'
    if(data.hasOwnProperty('cct')) return 'cct'
}