package go;

public class Func_int_int{

	Interface_call_int_int func;

	public Func_int_int(Interface_call_int_int func){
		this.func = func;
	}

	public int call(int x){
		return this.func.call(x);
	}
}
