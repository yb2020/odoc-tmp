/*
 * Created Date: August 26th 2021, 10:02:08 am
 * Author: zhoupengcheng
 * -----
 * Last Modified: March 2nd 2022, 11:11:12 am
 */
// @ts-nocheck
import UndoManager from 'undo-manager';
const undo = new UndoManager();

window.undoManager = undo;

export default undo;
